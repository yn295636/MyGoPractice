package main

import (
	"context"
	"github.com/agiledragon/gomonkey"
	"github.com/gin-gonic/gin/json"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/stretchr/testify/assert"
	pb "github.com/yn295636/MyGoPractice/proto"
	"reflect"
	"testing"
)

func TestGreeterService(t *testing.T) {
	t.Run("SayHello", func(tt *testing.T) {
		asserting := assert.New(tt)
		req := &pb.HelloRequest{
			Name: "tester",
		}
		s := server{}
		resp, err := s.SayHello(context.Background(), req)
		asserting.NoError(err)
		asserting.Equal("Hello tester", resp.Message)
	})

	t.Run("StoreInMongo", func(tt *testing.T) {
		mongoClient, _ = InitMongoClient()
		asserting := assert.New(tt)
		collection := []interface{}{}
		var mongoCollection *mongo.Collection
		patches := gomonkey.ApplyMethod(reflect.TypeOf(mongoCollection),
			"InsertOne", func(_ *mongo.Collection, _ context.Context, document interface{},
				opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
				collection = append(collection, document)
				return nil, nil
			})
		defer patches.Reset()

		jsonStr := `{"name": "tester"}`
		s := server{}
		resp, err := s.StoreInMongo(context.Background(), &pb.StoreInMongoRequest{
			Data: jsonStr,
		})
		asserting.NoError(err)
		asserting.Equal(int32(1), resp.Result)
		asserting.Equal(1, len(collection))
		if len(collection) == 1 {
			actualJsonBytes, err := json.Marshal(collection[0])
			asserting.NoError(err)
			asserting.JSONEq(jsonStr, string(actualJsonBytes))
		}
	})
}
