package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/yn295636/MyGoPractice/proto/greeter_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

// server is used to implement greeter_service.GreeterServer.
type server struct{}

// SayHello implements greeter_service.GreeterServer
func (s *server) SayHello(ctx context.Context, in *greeter_service.HelloRequest) (*greeter_service.HelloReply, error) {
	log.Printf("SayHello received: %v", in.Name)
	return &greeter_service.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) StoreInMongo(ctx context.Context, in *greeter_service.StoreInMongoRequest) (*greeter_service.StoreInMongoReply, error) {
	log.Printf("StoreInMongo received: %v", in.Data)
	var out *greeter_service.StoreInMongoReply
	out = &greeter_service.StoreInMongoReply{
		Result: 0,
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(in.Data), &data); err != nil {
		log.Printf("Json unmarshal error, %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := mongoClient.Database(MyMongoDb).Collection(MyMongoCollection).InsertOne(context.Background(), data); err != nil {
		log.Printf("Insert data into mongo failed, %v", err)
		return out, status.Error(codes.Internal, err.Error())
	}
	out.Result = 1
	return out, nil
}

func (s *server) StoreInRedis(ctx context.Context, in *greeter_service.StoreInRedisRequest) (*greeter_service.StoreInRedisReply, error) {
	log.Printf("StoreInRedis received: %v", in)
	var out *greeter_service.StoreInRedisReply
	out = &greeter_service.StoreInRedisReply{
		Result: 0,
	}
	redisConn := redisPool.Get()
	if _, err := redisConn.Do("SET", fmt.Sprintf("%v_%v", RedisPrefix, in.Key), in.Value); err != nil {
		log.Printf("Insert data into redis failed, %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	out.Result = 1
	return out, nil
}

func (s *server) GetFromRedis(ctx context.Context, in *greeter_service.GetFromRedisRequest) (*greeter_service.GetFromRedisReply, error) {
	log.Printf("GetFromRedis received: %v", in)
	var out *greeter_service.GetFromRedisReply
	redisConn := redisPool.Get()
	if result, err := redis.String(
		redisConn.Do("GET", fmt.Sprintf("%v_%v", RedisPrefix, in.Key)));
		err != nil {
		log.Printf("Get data from redis failed, %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		log.Printf("Got data from redis: %v", result)
		out = &greeter_service.GetFromRedisReply{Value: result}
		return out, nil
	}
}

func (s *server) StoreUserInDb(ctx context.Context, in *greeter_service.StoreUserInDbRequest) (*greeter_service.StoreUserInDbReply, error) {
	log.Printf("StoreUserInDb received: %v", in)
	var out *greeter_service.StoreUserInDbReply
	out = &greeter_service.StoreUserInDbReply{}
	id, err := StoreUser(in)
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		log.Printf("username already exists")
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	} else if err != nil {
		log.Printf("Insert user into db failed, %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	out.Uid = id
	return out, nil
}

func (s *server) StorePhoneInDb(ctx context.Context, in *greeter_service.StorePhoneInDbRequest) (*greeter_service.StorePhoneInDbReply, error) {
	log.Printf("StorePhoneInDb received: %v", in)
	var out *greeter_service.StorePhoneInDbReply
	id, err := StorePhone(in)
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		log.Printf("phone already exists")
		return nil, status.Error(codes.AlreadyExists, "phone already exists")
	} else if err != nil {
		log.Printf("Insert phone into db failed, %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	out = &greeter_service.StorePhoneInDbReply{
		Id: id,
	}
	return out, nil
}

func (s *server) GetUserFromDb(ctx context.Context, in *greeter_service.GetUserFromDbRequest) (*greeter_service.GetUserFromDbReply, error) {
	log.Printf("GetUserFromDb received: %v", in)
	var (
		out *greeter_service.GetUserFromDbReply
		err error
	)
	out, err = QueryUserByUid(in)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return nil, status.Error(codes.NotFound, "user doesn't exist")
	} else if err != nil {
		log.Printf("Query user by uid failed, %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return out, nil
}
