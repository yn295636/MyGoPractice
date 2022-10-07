package main

import (
	"fmt"
	"github.com/yn295636/MyGoPractice/app/apigateway/redis"
	"github.com/yn295636/MyGoPractice/app/apigateway/router"
	"github.com/yn295636/MyGoPractice/grpcfactory"
)

func main() {
	LoadConfig()
	grpcfactory.SetupGrpcClientFactory(GetSettings().EtcdAddrs)

	redis.InitRedisPool(GetSettings().RedisAddr)

	r := router.NewRouter()
	r.Run(fmt.Sprintf("0.0.0.0:%v", GetSettings().ListenPort))
}
