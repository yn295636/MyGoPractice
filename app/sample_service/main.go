package main

import (
	"github.com/yn295636/MyGoPractice/proto/sample_service"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	Port = ":50052"
)

func main() {
	var err error
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	sample_service.RegisterSampleServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
