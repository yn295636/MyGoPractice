/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	pb "github.com/yn295636/MyGoPractice/proto"
	"google.golang.org/grpc"
)

const (
	Port              = ":50051"
	MyMongoDb         = "mygopracticedb"
	MyMongoCollection = "mygopracticecollection"
	RedisPrefix       = "mygopractice"
	MongoAddr = ":27017"
	RedisAddr = "127.0.0.1:6379"
)

var (
	mongoClient *mongo.Client
	redisPool   *redis.Pool
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) StoreInMongo(ctx context.Context, in *pb.StoreInMongoRequest) (*pb.StoreInMongoReply, error) {
	log.Printf("Received: %v", in.Data)
	var out *pb.StoreInMongoReply
	out = &pb.StoreInMongoReply{
		Result: 0,
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(in.Data), &data); err != nil {
		log.Printf("Json unmarshal error, %v", err)
		return out, err
	}

	if _, err := mongoClient.Database(MyMongoDb).Collection(MyMongoCollection).InsertOne(context.Background(), data); err != nil {
		log.Printf("Insert data into mongo failed, %v", err)
		return out, err
	}
	out.Result = 1
	return out, nil
}

func (s *server) StoreInRedis(ctx context.Context, in *pb.StoreInRedisRequest) (*pb.StoreInRedisReply, error) {
	var out *pb.StoreInRedisReply
	out = &pb.StoreInRedisReply{
		Result: 0,
	}
	redisConn := redisPool.Get()
	if _, err := redisConn.Do("SET", fmt.Sprintf("%v_%v", RedisPrefix, in.Key), in.Value);
		err != nil {
		log.Printf("Insert data into redis failed, %v", err)
		return out, err
	}
	out.Result = 1
	return out, nil
}

func main() {
	mongoClient, _ = InitMongoClient(MongoAddr)
	redisPool = InitRedisPool(RedisAddr)

	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
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
		log.Fatalf("Init mongo client error, %v", err)
		return nil, err
	}
	if err = clt.Connect(context.TODO()); err != nil {
		log.Fatalf("Connect mongo db error, %v", err)
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
