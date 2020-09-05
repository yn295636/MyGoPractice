package grpcfactory

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/yn295636/MyGoPractice/etcd"
	greeterpb "github.com/yn295636/MyGoPractice/proto/greeter_service"
	samplepb "github.com/yn295636/MyGoPractice/proto/sample_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const (
	//GreeterAddr    = "localhost:50051"
	//SampleAddr     = "localhost:50052"
	SampleService  = "sample_service"
	GreeterService = "greeter_service"
)

var (
	grpcCF grpcClientFactoryI
)

type releaseFunc func()

type grpcClientFactoryI interface {
	newSampleClient() (samplepb.SampleClient, error, releaseFunc)
	newGreeterClient() (greeterpb.GreeterClient, error, releaseFunc)
}

func SetupGrpcClientFactory(etcdAddrs []string) {
	grpcCF = newGrpcClientFactory(etcdAddrs)
}

func NewSampleClient() (samplepb.SampleClient, error, releaseFunc) {
	return grpcCF.newSampleClient()
}

func NewGreeterClient() (greeterpb.GreeterClient, error, releaseFunc) {
	return grpcCF.newGreeterClient()
}

type grpcClientFactory struct {
	EtcdAddrs []string
}

func newGrpcClientFactory(etcdAddrs []string) grpcClientFactoryI {
	return &grpcClientFactory{
		EtcdAddrs: etcdAddrs,
	}
}

func (f *grpcClientFactory) newSampleClient() (samplepb.SampleClient, error, releaseFunc) {
	r := etcd.NewResolver(f.EtcdAddrs, SampleService)
	resolver.Register(r)
	addr := fmt.Sprintf("%s:///%s", r.Scheme(), SampleService)
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithBalancerName("round_robin"),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true}))
	//conn, err := grpc.Dial(SampleAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err, func() {}
	}
	c := samplepb.NewSampleClient(conn)
	return c, nil, func() {
		err = conn.Close()
		if err != nil {
			log.Printf("release connection error: %v", err)
		}
	}
}

func (f *grpcClientFactory) newGreeterClient() (greeterpb.GreeterClient, error, releaseFunc) {
	r := etcd.NewResolver(f.EtcdAddrs, GreeterService)
	resolver.Register(r)
	addr := fmt.Sprintf("%s:///%s", r.Scheme(), GreeterService)
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithBalancerName("round_robin"),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true}))
	//conn, err := grpc.Dial(GreeterAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err, func() {}
	}
	c := greeterpb.NewGreeterClient(conn)
	return c, nil, func() {
		err = conn.Close()
		if err != nil {
			log.Printf("release connection error: %v", err)
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
