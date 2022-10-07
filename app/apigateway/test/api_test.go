package test

import (
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	gofight "github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"github.com/yn295636/MyGoPractice/app/apigateway/redis"
	"github.com/yn295636/MyGoPractice/app/apigateway/router"
	"github.com/yn295636/MyGoPractice/grpcfactory"
	pb "github.com/yn295636/MyGoPractice/proto/greeter_service"
)

var (
	_, debug        = os.LookupEnv("DEBUG")
	ginEng          *gin.Engine
	mockGrpcFactory *grpcfactory.MockGrpcClientFactory
	testClient      *gofight.RequestConfig
	mredis          *miniredis.Miniredis
)

func init() {
	gin.SetMode(gin.TestMode)
	ginEng = router.NewRouter()
	testClient = gofight.New()
}

func Setup(t *testing.T) func() {
	asserting := assert.New(t)
	ctrl := gomock.NewController(t)
	mockGrpcFactory = grpcfactory.SetupMockGrpcClientFactory(ctrl)

	var err error
	mredis, err = miniredis.Run()
	asserting.NoError(err)
	redis.InitRedisPool(mredis.Addr())

	return func() {
		mockGrpcFactory = nil
		ctrl.Finish()
		mredis.Close()
	}
}

func TestGreeting(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown()
	mockHelloReq := &pb.HelloRequest{
		Name: "tester",
	}
	mockGrpcFactory.GreeterClient.EXPECT().
		SayHello(gomock.Any(), gomock.Eq(mockHelloReq)).
		Return(
			&pb.HelloReply{
				Message: "Hello tester",
			}, nil).
		MinTimes(1)
	t.Run("Success", func(tt *testing.T) {
		defer mredis.FlushAll()
		asserting := assert.New(tt)
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

	t.Run("RateLimit", func(tt *testing.T) {
		defer mredis.FlushAll()
		asserting := assert.New(tt)
		wg := sync.WaitGroup{}
		total := 15
		threshold := 5
		successCount := int32(0)
		failedCount := int32(0)
		for i := 0; i < total; i++ {
			time.Sleep(time.Millisecond * 10)
			wg.Add(1)
			go func() {
				defer wg.Done()
				testClient.POST("/greet").
					SetJSON(map[string]interface{}{
						"name": "tester",
					}).
					SetDebug(debug).
					Run(ginEng, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
						if resp.Code == http.StatusTooManyRequests {
							atomic.AddInt32(&failedCount, 1)
						} else if resp.Code == http.StatusOK {
							atomic.AddInt32(&successCount, 1)
						}
					})
			}()
		}
		wg.Wait()
		asserting.EqualValues(successCount, threshold)
		asserting.EqualValues(successCount+failedCount, total)
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
