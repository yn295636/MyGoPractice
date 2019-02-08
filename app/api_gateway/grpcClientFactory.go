package main

import (
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

func NewMockGrpcClientFactory() GrpcClientFactory {
	return &mockGrpcClientFactory{}
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

type mockGrpcClientFactory struct {
	greeterClient pb.GreeterClient
}

func (f *mockGrpcClientFactory) NewGreeterClient() (pb.GreeterClient, error, releaseFunc) {
	return f.greeterClient, nil, func() {}
}
