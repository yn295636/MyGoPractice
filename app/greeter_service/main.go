package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/yn295636/MyGoPractice/db"
	"log"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yn295636/MyGoPractice/proto/greeter_service"
	"google.golang.org/grpc"
)

const (
	Port              = ":50051"
	MyMongoDb         = "mygopracticedb"
	MyMongoCollection = "mygopracticecollection"
	RedisPrefix       = "mygopractice"
	MongoAddr         = ":27017"
	RedisAddr         = "127.0.0.1:6379"
	DbAddr            = "127.0.0.1:3306"
	DbUser            = "root"
	DbPassword        = "Mygopractice123!"
)

var (
	mongoClient *mongo.Client
	redisPool   *redis.Pool
)

func main() {
	var err error

	mongoClient, err = InitMongoClient(MongoAddr)
	if err != nil {
		log.Fatalf("Init Mongo failed, %v", err)
	}

	redisPool = InitRedisPool(RedisAddr)

	err = db.InitDb(DbAddr, DbUser, DbPassword, db.DbLatestVer)
	if err != nil {
		log.Fatalf("Init Db failed, %v", err)
	}

	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greeter_service.RegisterGreeterServer(s, &server{})
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
