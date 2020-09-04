package test

import (
	gofight "github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"github.com/yn295636/MyGoPractice/app/apigateway/router"
	"github.com/yn295636/MyGoPractice/grpcfactory"
	pb "github.com/yn295636/MyGoPractice/proto/greeter_service"
	"net/http"
	"os"
	"testing"
)

var (
	_, debug        = os.LookupEnv("DEBUG")
	ginEng          *gin.Engine
	mockGrpcFactory *grpcfactory.MockGrpcClientFactory
	testClient      *gofight.RequestConfig
)

func init() {
	gin.SetMode(gin.TestMode)
	ginEng = router.NewRouter()
	testClient = gofight.New()
}

func Setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockGrpcFactory = grpcfactory.SetupMockGrpcClientFactory(ctrl)

	return func() {
		mockGrpcFactory = nil
		ctrl.Finish()
	}
}

func TestGreeting(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown()
	asserting := assert.New(t)
	mockHelloReq := &pb.HelloRequest{
		Name: "tester",
	}
	mockGrpcFactory.GreeterClient.EXPECT().
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
}

func TestStoreUserInDb(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown()
	asserting := assert.New(t)
	var (
		username    = "yn"
		gender      = int32(2)
		age         = int32(20)
		uid         = int64(10)
		externalUid = int32(123)
	)

	mockStoreUserInDbReq := &pb.StoreUserInDbRequest{
		Username:    username,
		Gender:      gender,
		Age:         age,
		ExternalUid: externalUid,
	}
	mockStoreUserInDbReply := &pb.StoreUserInDbReply{
		Uid: uid,
	}
	mockGrpcFactory.GreeterClient.EXPECT().
		StoreUserInDb(gomock.Any(), gomock.Eq(mockStoreUserInDbReq)).
		Return(mockStoreUserInDbReply, nil)
	testClient.POST("/db/users").
		SetJSON(map[string]interface{}{
			"username":     username,
			"gender":       gender,
			"age":          age,
			"external_uid": externalUid,
		}).
		SetDebug(debug).
		Run(ginEng, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
			expected := `
                    {
                      "uid": 10
                    }
                `
			asserting.Equal(http.StatusOK, resp.Code)
			asserting.JSONEq(expected, resp.Body.String())
		})
}
