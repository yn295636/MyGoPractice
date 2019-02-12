package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	pb "github.com/yn295636/MyGoPractice/proto"
	"log"
	"net/http"
	"time"
)

func Greet(c *gin.Context) {
	var body GreetReq
	if err := c.ShouldBind(&body); err != nil {
		log.Printf("Failed to bind to GreetReq, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}

	client, err, release := grpcCF.NewGreeterClient()
	if err != nil {
		log.Printf("Failed to get greeter client, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		return
	}
	defer release()

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		time.Duration(2*time.Second))
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: body.Name,
	})
	if err != nil {
		log.Printf("Failed to say hello, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		return
	}
	c.JSON(http.StatusOK, GreetResp{
		Message: resp.Message,
	})
}

func StoreInMongo(c *gin.Context) {
	var body StoreInMongoReq
	if err := c.ShouldBind(&body); err != nil {
		log.Printf("Failed to bind to StoreInMongoReq, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to json marshal request body, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}

	client, err, release := grpcCF.NewGreeterClient()
	if err != nil {
		log.Printf("Failed to get greeter client, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		return
	}
	defer release()

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		time.Duration(2*time.Second))
	defer cancel()

	resp, err := client.StoreInMongo(ctx, &pb.StoreInMongoRequest{
		Data: string(jsonBytes),
	})
	if err != nil {
		log.Printf("Failed to store in mongo, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		return
	}
	c.JSON(http.StatusOK, StoreInMongoResp{
		Result: resp.Result,
	})
}

func StoreInRedis(c *gin.Context) {
	var body StoreInRedisReq
	if err := c.ShouldBind(&body); err != nil {
		log.Printf("Failed to bind to StoreInRedisReq, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}

	client, err, release := grpcCF.NewGreeterClient()
	if err != nil {
		log.Printf("Failed to get greeter client, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		return
	}
	defer release()

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		time.Duration(2*time.Second))
	defer cancel()

	resp, err := client.StoreInRedis(ctx, &pb.StoreInRedisRequest{
		Key:   body.Key,
		Value: body.Value,
	})
	if err != nil {
		log.Printf("Failed to store in redis, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		return
	}
	c.JSON(http.StatusOK, StoreInRedisResp{
		Result: resp.Result,
	})
}
