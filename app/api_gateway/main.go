package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yn295636/MyGoPractice/db"
	"log"
)

const (
	Port       = 8080
	DbAddr     = "127.0.0.1:3306"
	DbUser     = "tester"
	DbPassword = "tester"
)

var (
	grpcCF GrpcClientFactory
)

func main() {
	grpcCF = NewGrpcClientFactory()

	err := db.InitDb(DbAddr, DbUser, DbPassword, db.DbVer2)
	if err != nil {
		log.Fatalf("Init Db failed, %v", err)
	}

	r := Router()
	r.Run(fmt.Sprintf("0.0.0.0:%v", Port))
}

func Router() *gin.Engine {
	r := gin.Default()
	r.POST("/greet", Greet)
	r.POST("/storeinmongo", StoreInMongo)
	r.POST("/storeinredis", StoreInRedis)
	return r
}
