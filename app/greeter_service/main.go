package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/yn295636/MyGoPractice/common"
	"github.com/yn295636/MyGoPractice/db"
	"github.com/yn295636/MyGoPractice/etcd"
	"github.com/yn295636/MyGoPractice/grpcfactory"
	"log"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yn295636/MyGoPractice/proto/greeter_service"
	"google.golang.org/grpc"
)

const (
	MyMongoDb         = "mygopracticedb"
	MyMongoCollection = "mygopracticecollection"
	RedisPrefix       = "mygopractice"
)

var (
	mongoClient *mongo.Client
	redisPool   *redis.Pool
)

func main() {
	LoadConfig()
	var err error

	mongoClient, err = InitMongoClient(GetSettings().MongoAddr)
	if err != nil {
		log.Fatalf("Init Mongo failed, %v", err)
	}

	redisPool = InitRedisPool(GetSettings().RedisAddr)

	err = db.InitDb(GetSettings().DbAddr, GetSettings().DbUser, GetSettings().DbPassword, db.DbLatestVer)
	if err != nil {
		log.Fatalf("Init Db failed, %v", err)
	}

	grpcfactory.SetupGrpcClientFactory(GetSettings().EtcdAddrs)

	lis, err := net.Listen("tcp", GetSettings().ListenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greeter_service.RegisterGreeterServer(s, &server{})

	// Register on etcd
	reg, err := etcd.NewService(etcd.ServiceInfo{
		Name: "greeter_service",
		IP:   fmt.Sprintf("%v%v", common.GetIPAddr(), GetSettings().ListenPort),
	}, GetSettings().EtcdAddrs)
	if err != nil {
		log.Fatal(err)
	}
	go reg.Start()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func InitMongoClient(mongoAddr string) (*mongo.Client, error) {
	var clt *mongo.Client
	var err error
	if clt, err = mongo.NewClient(fmt.Sprintf(
		"mongodb://%s:%s@%s",
		"",
		"",
		mongoAddr)); err != nil {
		log.Printf("Init mongo client error, %v", err)
		return nil, err
	}
	if err = clt.Connect(context.TODO()); err != nil {
		log.Printf("Connect mongo db error, %v", err)
		return nil, err
	}
	return clt, nil
}

func InitRedisPool(redisAddr string) *redis.Pool {
	var pool *redis.Pool
	pool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}
