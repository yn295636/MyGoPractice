package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yn295636/MyGoPractice/app/apigateway/api"
	"github.com/yn295636/MyGoPractice/app/apigateway/sample_api"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/greet", api.Greet)
	r.POST("/mongo", api.StoreInMongo)
	r.POST("/redis", api.StoreInRedis)
	r.POST("/db/users", api.StoreUserInDb)
	r.POST("/db/users/:uid/phones", api.StoreUserPhoneInDb)
	r.GET("/multiply", sample_api.Multiple)
	return r
}
