package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yn295636/MyGoPractice/grpcfactory"
	pb "github.com/yn295636/MyGoPractice/proto/greeter_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"strconv"
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

	client, err, release := grpcfactory.NewGreeterClient()
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

	client, err, release := grpcfactory.NewGreeterClient()
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
			Message: err.Error(),
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

	client, err, release := grpcfactory.NewGreeterClient()
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
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, StoreInRedisResp{
		Result: resp.Result,
	})
}

func StoreUserInDb(c *gin.Context) {
	var body StoreUserInDbReq
	if err := c.ShouldBind(&body); err != nil {
		log.Printf("Failed to bind to StoreUserInDbReq, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}

	client, err, release := grpcfactory.NewGreeterClient()
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
		2*time.Second)
	defer cancel()

	resp, err := client.StoreUserInDb(ctx, &pb.StoreUserInDbRequest{
		Username:    body.Username,
		Gender:      int32(body.Gender),
		Age:         int32(body.Age),
		ExternalUid: body.ExternalUid,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			c.JSON(http.StatusBadRequest, ErrorRsp{
				Code:    http.StatusBadRequest,
				Message: s.Message(),
			})
			return
		}
		log.Printf("Failed to store user in db, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, StoreUserInDbResp{
		Uid: resp.Uid,
	})
}

func StoreUserPhoneInDb(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		log.Printf("Failed to parse uid from path, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}

	var body StoreUserPhoneInDbReq
	if err := c.ShouldBind(&body); err != nil {
		log.Printf("Failed to bind to StoreUserPhoneInDbReq, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}

	client, err, release := grpcfactory.NewGreeterClient()
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

	_, err = client.GetUserFromDb(ctx, &pb.GetUserFromDbRequest{Uid: uid})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			c.JSON(http.StatusBadRequest, ErrorRsp{
				Code:    http.StatusBadRequest,
				Message: s.Message(),
			})
			return
		}

		log.Printf("Failed to get user from db, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	_, err = client.StorePhoneInDb(ctx, &pb.StorePhoneInDbRequest{
		Uid:   uid,
		Phone: body.Phone,
	})
	if err != nil {
		log.Printf("Failed to store phone into db, %v", err)
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			c.JSON(http.StatusBadRequest, ErrorRsp{
				Code:    http.StatusBadRequest,
				Message: s.Message(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, StoreUserPhoneInDbResp{})
}
