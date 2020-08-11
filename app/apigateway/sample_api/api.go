package sample_api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yn295636/MyGoPractice/app/apigateway/grpcfactory"
	pb "github.com/yn295636/MyGoPractice/proto/sample_service"
	"log"
	"net/http"
	"time"
)

func Multiple(c *gin.Context) {
	var body MultiplyReq
	if err := c.ShouldBind(&body); err != nil {
		log.Printf("Failed to bind to MultiplyReq, %v", err)
		c.JSON(http.StatusBadRequest, ErrorRsp{
			Code:    http.StatusBadRequest,
			Message: "input error",
		})
		return
	}

	sampleClient, err, release := grpcfactory.NewSampleClient()
	if err != nil {
		log.Printf("Failed to get sample sample_client, %v", err)
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

	resp, err := sampleClient.Multiply(ctx, &pb.MultiplyReq{
		A: body.A,
		B: body.B,
	})

	if err != nil {
		log.Printf("Failed to invoke sample_service.Multiply, %v", err)
		c.JSON(http.StatusInternalServerError, ErrorRsp{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		return
	}
	c.JSON(http.StatusOK, MultiplyResp{
		Result: resp.Result,
	})
}
