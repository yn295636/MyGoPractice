package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc/grpc-go/status"
	pb "github.com/yn295636/MyGoPractice/proto"
	"google.golang.org/grpc/codes"
	"log"
	"strings"
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
		return nil, err
	}

	if _, err := mongoClient.Database(MyMongoDb).Collection(MyMongoCollection).InsertOne(context.Background(), data); err != nil {
		log.Printf("Insert data into mongo failed, %v", err)
		return out, err
	}
	out.Result = 1
	return out, nil
}

func (s *server) StoreInRedis(ctx context.Context, in *pb.StoreInRedisRequest) (*pb.StoreInRedisReply, error) {
	log.Printf("Received: %v", in)
	var out *pb.StoreInRedisReply
	out = &pb.StoreInRedisReply{
		Result: 0,
	}
	redisConn := redisPool.Get()
	if _, err := redisConn.Do("SET", fmt.Sprintf("%v_%v", RedisPrefix, in.Key), in.Value);
		err != nil {
		log.Printf("Insert data into redis failed, %v", err)
		return nil, err
	}
	out.Result = 1
	return out, nil
}

func (s *server) StoreUserInDb(ctx context.Context, in *pb.StoreUserInDbRequest) (*pb.StoreUserInDbReply, error) {
	log.Printf("Received: %v", in)
	var out *pb.StoreUserInDbReply
	out = &pb.StoreUserInDbReply{}
	id, err := StoreUser(in)
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		log.Printf("username already exists")
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}
	out.Uid = id
	return out, nil
}

func (s *server) StorePhoneInDb(ctx context.Context, in *pb.StorePhoneInDbRequest) (*pb.StorePhoneInDbReply, error) {
	log.Printf("Received: %v", in)
	var out *pb.StorePhoneInDbReply
	id, err := StorePhone(in)
	if err != nil {
		return nil, err
	}
	out = &pb.StorePhoneInDbReply{
		Id: id,
	}
	return out, nil
}
