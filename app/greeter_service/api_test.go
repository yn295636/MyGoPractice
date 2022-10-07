package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/stretchr/testify/require"
	"github.com/yn295636/MyGoPractice/db"
	pb "github.com/yn295636/MyGoPractice/proto/greeter_service"
)

var (
	testDbAddr     = "127.0.0.1:13306"
	testDbUser     = "root"
	testDbPassword = "Mytest123!"
)

func TestSayHello(t *testing.T) {
	asserting := require.New(t)
	req := &pb.HelloRequest{
		Name: "tester",
	}
	s := server{}
	resp, err := s.SayHello(context.Background(), req)
	asserting.NoError(err)
	asserting.Equal("Hello tester", resp.Message)
}

func TestStoreInMongo(t *testing.T) {
	mongoClient, _ = InitMongoClient(GetSettings().MongoAddr)
	asserting := require.New(t)
	var collection []interface{}
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
}

func TestStoreInRedis(t *testing.T) {
	asserting := require.New(t)
	mredis, err := miniredis.Run()
	asserting.NoError(err)
	defer mredis.Close()
	redisPool = InitRedisPool(mredis.Addr())

	key := "name"
	value := "tester"

	s := server{}
	resp, err := s.StoreInRedis(context.Background(), &pb.StoreInRedisRequest{
		Key:   key,
		Value: value,
	})
	asserting.NoError(err)
	asserting.Equal(int32(1), resp.Result)
	actualVal, err := mredis.Get(fmt.Sprintf("%v_%v", RedisPrefix, key))
	asserting.NoError(err)
	asserting.Equal(value, actualVal)
}

func TestStoreUserInDb_WithSuccess(t *testing.T) {
	tearDownDb := initDbForTest(t)
	defer tearDownDb()

	assert := require.New(t)

	defer clearUserTable(t)
	s := server{}
	req := &pb.StoreUserInDbRequest{
		Username:    "test user",
		Gender:      1,
		Age:         20,
		ExternalUid: 123,
	}
	resp, err := s.StoreUserInDb(context.Background(), req)
	assert.NoError(err)

	result, err := QueryUserByUid(&pb.GetUserFromDbRequest{
		Uid: resp.Uid,
	})
	assert.NoError(err)
	assert.Equal("test user", result.Username)
	assert.EqualValues(1, result.Gender)
	assert.EqualValues(20, result.Age)
	assert.EqualValues(123, result.ExternalUid)
}

func TestStoreUserInDb_WithDup(t *testing.T) {
	tearDownDb := initDbForTest(t)
	defer tearDownDb()

	assert := require.New(t)

	defer clearUserTable(t)

	s := server{}
	req := &pb.StoreUserInDbRequest{
		Username:    "test user2",
		Gender:      2,
		Age:         22,
		ExternalUid: 122,
	}
	_, err := s.StoreUserInDb(context.Background(), req)
	assert.NoError(err)
	_, err = s.StoreUserInDb(context.Background(), req)
	assert.Error(err)
	assert.Contains(err.Error(), "username already exists")
}

func clearUserTable(t *testing.T) {
	assert := require.New(t)
	dbConn := db.GetDb(db.DbMyGoPractice)
	sqlStat := fmt.Sprintf("DELETE FROM %v;", db.TblUser)
	_, err := dbConn.Exec(sqlStat)
	assert.NoError(err)
}

func initDbForTest(t *testing.T) func() {
	assert := require.New(t)
	err := db.InitDb(testDbAddr, testDbUser, testDbPassword, db.DbLatestVer)
	assert.NoError(err)
	return func() {
		db.DisconnectAllDb()
	}
}
