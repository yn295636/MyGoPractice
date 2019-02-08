package main

import (
	"github.com/golang/mock/gomock"
	pb "github.com/yn295636/MyGoPractice/proto"
	"google.golang.org/grpc"
	"log"
)

const (
	GreeterAddr = "localhost:50051"
)

type releaseFunc func()

type GrpcClientFactory interface {
	NewGreeterClient() (pb.GreeterClient, error, releaseFunc)
}

func NewGrpcClientFactory() GrpcClientFactory {
	return &grpcClientFactory{}
}

func NewMockGrpcClientFactory(ctrl *gomock.Controller) GrpcClientFactory {
	return &MockGrpcClientFactory{
		greeterClient: pb.NewMockGreeterClient(ctrl),
	}
}

type grpcClientFactory struct {
}

func (f *grpcClientFactory) NewGreeterClient() (pb.GreeterClient, error, releaseFunc) {
	conn, err := grpc.Dial(GreeterAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err, func() {}
	}
	c := pb.NewGreeterClient(conn)
	return c, nil, func() {
		err = conn.Close()
		if err != nil {
			log.Printf("release connection error: %v", err)
		}
	}
}

type MockGrpcClientFactory struct {
	greeterClient *pb.MockGreeterClient
}

func (f *MockGrpcClientFactory) NewGreeterClient() (pb.GreeterClient, error, releaseFunc) {
	return f.greeterClient, nil, func() {}
}
