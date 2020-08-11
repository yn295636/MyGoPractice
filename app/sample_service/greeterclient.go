package main

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/yn295636/MyGoPractice/proto/greeter_service"
	"google.golang.org/grpc"
	"log"
)

const (
	GreeterServicePort = 50051
)

var (
	isTest            bool = false
	mockGreeterClient greeter_service.GreeterClient
)

func GetGreeterClient() (greeter_service.GreeterClient, error, func()) {
	if !isTest {
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", GreeterServicePort), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
			return nil, err, func() {}
		}
		c := greeter_service.NewGreeterClient(conn)
		return c, nil, func() {
			err = conn.Close()
			if err != nil {
				log.Fatalf("release connection error: %v", err)
			}
		}
	} else {
		if mockGreeterClient == nil {
			return nil, errors.New("mockGreeterClient is not initialized"), func() {}
		}
		return mockGreeterClient, nil, func() {
			mockGreeterClient = nil
		}
	}
}

func NewMockGreeterClient(ctrl *gomock.Controller) greeter_service.GreeterClient {
	isTest = true
	mockGreeterClient = greeter_service.NewMockGreeterClient(ctrl)
	return mockGreeterClient
}
