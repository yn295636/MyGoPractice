package main

import (
	"fmt"

	"github.com/yn295636/MyGoPractice/app/apigateway/grpcfactory"
	"github.com/yn295636/MyGoPractice/app/apigateway/router"
)

const (
	Port = 8081
)

func main() {
	grpcfactory.SetupGrpcClientFactory()

	r := router.NewRouter()
	r.Run(fmt.Sprintf("0.0.0.0:%v", Port))
}
