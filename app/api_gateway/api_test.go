package main

import (
	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	pb "github.com/yn295636/MyGoPractice/proto"
	"net/http"
	"os"
	"testing"
)

var (
	_, debug = os.LookupEnv("DEBUG")
)

func Setup(t *testing.T, ctrl *gomock.Controller) func() {
	gin.SetMode(gin.TestMode)
	grpcCF = NewMockGrpcClientFactory(ctrl)
	return func() {

	}
}

func TestApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tearDown := Setup(t, ctrl)
	defer tearDown()
	ginEng := Router()
	testClient := gofight.New()
	t.Run("Greet", func(tt *testing.T) {
		asserting := assert.New(tt)
		mockHelloReq := &pb.HelloRequest{
			Name: "tester",
		}
		grpcCF.(*MockGrpcClientFactory).greeterClient.EXPECT().
			SayHello(gomock.Any(), gomock.Eq(mockHelloReq)).
			Return(
				&pb.HelloReply{
					Message: "Hello tester",
				}, nil)
		testClient.POST("/greet").
			SetJSON(map[string]interface{}{
				"name": "tester",
			}).
			SetDebug(debug).
			Run(ginEng, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
				expected := `
					{
					  "message": "Hello tester"
					}
				`
				asserting.Equal(http.StatusOK, resp.Code)
				asserting.JSONEq(expected, resp.Body.String())
			})
	})
}
