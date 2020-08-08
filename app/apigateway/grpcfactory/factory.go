package grpcfactory

import (
	"github.com/golang/mock/gomock"
	greeterpb "github.com/yn295636/MyGoPractice/proto/greeter_service"
	samplepb "github.com/yn295636/MyGoPractice/proto/sample_service"
	"google.golang.org/grpc"
	"log"
)

const (
	GreeterAddr = "localhost:50051"
	SampleAddr  = "localhost:50052"
)

var (
	grpcCF grpcClientFactoryI
)

type releaseFunc func()

type grpcClientFactoryI interface {
	newSampleClient() (samplepb.SampleClient, error, releaseFunc)
	newGreeterClient() (greeterpb.GreeterClient, error, releaseFunc)
}

func SetupGrpcClientFactory() {
	grpcCF = newGrpcClientFactory()
}

func NewSampleClient() (samplepb.SampleClient, error, releaseFunc) {
	return grpcCF.newSampleClient()
}

func NewGreeterClient() (greeterpb.GreeterClient, error, releaseFunc) {
	return grpcCF.newGreeterClient()
}

type grpcClientFactory struct {
}

func newGrpcClientFactory() grpcClientFactoryI {
	return &grpcClientFactory{}
}

func (f *grpcClientFactory) newSampleClient() (samplepb.SampleClient, error, releaseFunc) {
	conn, err := grpc.Dial(SampleAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err, func() {}
	}
	c := samplepb.NewSampleClient(conn)
	return c, nil, func() {
		err = conn.Close()
		if err != nil {
			log.Fatalf("release connection error: %v", err)
		}
	}
}

func (f *grpcClientFactory) newGreeterClient() (greeterpb.GreeterClient, error, releaseFunc) {
	conn, err := grpc.Dial(GreeterAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err, func() {}
	}
	c := greeterpb.NewGreeterClient(conn)
	return c, nil, func() {
		err = conn.Close()
		if err != nil {
			log.Fatalf("release connection error: %v", err)
		}
	}
}

type MockGrpcClientFactory struct {
	SampleClient  *samplepb.MockSampleClient
	GreeterClient *greeterpb.MockGreeterClient
}

func newMockGrpcClientFactory(ctrl *gomock.Controller) grpcClientFactoryI {
	return &MockGrpcClientFactory{
		SampleClient:  samplepb.NewMockSampleClient(ctrl),
		GreeterClient: greeterpb.NewMockGreeterClient(ctrl),
	}
}

func (f *MockGrpcClientFactory) newSampleClient() (samplepb.SampleClient, error, releaseFunc) {
	return f.SampleClient, nil, func() {}
}

func (f *MockGrpcClientFactory) newGreeterClient() (greeterpb.GreeterClient, error, releaseFunc) {
	return f.GreeterClient, nil, func() {}
}

func SetupMockGrpcClientFactory(ctrl *gomock.Controller) *MockGrpcClientFactory {
	mockClientFactory := newMockGrpcClientFactory(ctrl)
	grpcCF = mockClientFactory
	return mockClientFactory.(*MockGrpcClientFactory)
}
