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
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/yn295636/MyGoPractice/db"
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
	MongoAddr         = ":27017"
	RedisAddr         = "127.0.0.1:6379"
	DbAddr            = "127.0.0.1:3306"
	DbUser            = "tester"
	DbPassword        = "tester"
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
	pb.RegisterGreeterServer(s, &server{})
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
