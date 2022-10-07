package main

import (
	"fmt"
	"github.com/yn295636/MyGoPractice/common"
	"github.com/yn295636/MyGoPractice/etcd"
	"github.com/yn295636/MyGoPractice/grpcfactory"
	"github.com/yn295636/MyGoPractice/proto/sample_service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	LoadConfig()
	var err error

	grpcfactory.SetupGrpcClientFactory(GetSettings().EtcdAddrs)

	lis, err := net.Listen("tcp", GetSettings().ListenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	sample_service.RegisterSampleServer(s, &server{})

	// Register on etcd
	reg, err := etcd.NewService(etcd.ServiceInfo{
		Name: "sample_service",
		IP:   fmt.Sprintf("%v%v", common.GetIPAddr(), GetSettings().ListenPort),
	}, GetSettings().EtcdAddrs)
	if err != nil {
		log.Fatal(err)
	}
	go reg.Start()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
