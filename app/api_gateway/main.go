package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	Port = 8080
)

var (
	grpcCF GrpcClientFactory
)

func main() {
	grpcCF = NewGrpcClientFactory()

	r := Router()
	r.Run(fmt.Sprintf("0.0.0.0:%v", Port))
}

func Router() *gin.Engine {
	r := gin.Default()
	r.POST("/greet", Greet)
	r.POST("/mongo", StoreInMongo)
	r.POST("/redis", StoreInRedis)
	r.POST("/db/users", StoreUserInDb)
	r.POST("/db/users/:uid/phones", StoreUserPhoneInDb)
	return r
}
