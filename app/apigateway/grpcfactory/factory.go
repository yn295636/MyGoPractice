package grpcfactory

import (
	"github.com/golang/mock/gomock"
	pb "github.com/yn295636/MyGoPractice/proto/greeter_service"
	"google.golang.org/grpc"
	"log"
)

const (
	GreeterAddr = "localhost:50051"
)

var (
	grpcCF grpcClientFactoryI
)

type releaseFunc func()

type grpcClientFactoryI interface {
	newGreeterClient() (pb.GreeterClient, error, releaseFunc)
}

func SetupGrpcClientFactory() {
	grpcCF = newGrpcClientFactory()
}

func NewGreeterClient() (pb.GreeterClient, error, releaseFunc) {
	return grpcCF.newGreeterClient()
}

func SetupMockGrpcClientFactory(ctrl *gomock.Controller) *MockGrpcClientFactory {
	mockClientFactory := newMockGrpcClientFactory(ctrl)
	grpcCF = mockClientFactory
	return mockClientFactory.(*MockGrpcClientFactory)
}

func newGrpcClientFactory() grpcClientFactoryI {
	return &grpcClientFactory{}
}

type grpcClientFactory struct {
}

func (f *grpcClientFactory) newGreeterClient() (pb.GreeterClient, error, releaseFunc) {
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

func newMockGrpcClientFactory(ctrl *gomock.Controller) grpcClientFactoryI {
	return &MockGrpcClientFactory{
		GreeterClient: pb.NewMockGreeterClient(ctrl),
	}
}

type MockGrpcClientFactory struct {
	GreeterClient *pb.MockGreeterClient
}

func (f *MockGrpcClientFactory) newGreeterClient() (pb.GreeterClient, error, releaseFunc) {
	return f.GreeterClient, nil, func() {}
}
