package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dankinder/httpmock"
	"github.com/golang/mock/gomock"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yn295636/MyGoPractice/proto/greeter_service"
	"github.com/yn295636/MyGoPractice/proto/sample_service"
	"net/http"
	"testing"
	"time"
)

func TestMultiply(t *testing.T) {
	testDatas := []map[string]int64{
		{
			"a":      3,
			"b":      2,
			"result": 6,
		},
		{
			"a":      -3,
			"b":      2,
			"result": -6,
		},
	}
	assert := require.New(t)
	s := server{}
	for _, testData := range testDatas {
		a, _ := testData["a"]
		b, _ := testData["b"]
		result, _ := testData["result"]
		req := &sample_service.MultiplyReq{
			A: int32(a),
			B: int32(b),
		}
		resp, err := s.Multiply(context.Background(), req)
		assert.NoError(err)
		assert.EqualValues(result, resp.Result)
	}
}

func TestGetUserById_WithSuccess(t *testing.T) {
	externalUid := 2
	externalUserInfoPath := fmt.Sprintf("/users/%v", externalUid)
	externalUserInfoJson, _ := json.Marshal(map[string]interface{}{
		"name":   "external user",
		"gender": 1,
		"age":    18,
	})

	mockHttpHandler := &httpmock.MockHandler{}
	mockHttpHandler.
		On("Handle", "GET", externalUserInfoPath, mock.Anything).
		Return(httpmock.Response{
			Status: 200,
			Body:   externalUserInfoJson,
		})
	mockHttpServer := httpmock.NewServer(mockHttpHandler)
	defer mockHttpServer.Close()

	oriExternalHttpHost := externalHttpHost
	defer func() { externalHttpHost = oriExternalHttpHost }()
	externalHttpHost = mockHttpServer.URL()

	assert := require.New(t)
	s := server{}
	req := &sample_service.GetUserByIdReq{
		Uid: int32(externalUid),
	}
	resp, err := s.GetUserById(context.Background(), req)
	assert.NoError(err)
	assert.Equal("external user", resp.Name)
	assert.EqualValues(1, resp.Gender)
	assert.EqualValues(18, resp.Age)
}

func TestGetUserById_WithSuccess2(t *testing.T) {
	externalUid := 2
	externalUserInfoPath := fmt.Sprintf("/users/%v", externalUid)
	externalUserInfo := map[string]interface{}{
		"name":   "external user2",
		"gender": 2,
		"age":    20,
	}

	defer func() {
		gock.OffAll()
		gock.DefaultMatcher = gock.NewMatcher()
	}()
	gock.New(externalHttpHost).
		Get(externalUserInfoPath).
		Response.
		Status(http.StatusOK).
		JSON(externalUserInfo)

	assert := require.New(t)
	s := server{}
	req := &sample_service.GetUserByIdReq{
		Uid: int32(externalUid),
	}
	resp, err := s.GetUserById(context.Background(), req)
	assert.NoError(err)
	assert.Equal("external user2", resp.Name)
	assert.EqualValues(2, resp.Gender)
	assert.EqualValues(20, resp.Age)
}

func TestGetUserById_ExternalTimeout(t *testing.T) {
	externalUid := 2
	externalUserInfoPath := fmt.Sprintf("/users/%v", externalUid)
	externalUserInfoJson, _ := json.Marshal(map[string]interface{}{
		"name":   "external user",
		"gender": 1,
		"age":    18,
	})

	mockHttpHandler := &httpmock.MockHandler{}
	mockHttpHandler.
		On("Handle", "GET", externalUserInfoPath, mock.Anything).
		After(1100 * time.Millisecond).
		Return(httpmock.Response{
			Status: 200,
			Body:   externalUserInfoJson,
		})
	mockHttpServer := httpmock.NewServer(mockHttpHandler)
	defer mockHttpServer.Close()

	oriExternalHttpHost := externalHttpHost
	defer func() { externalHttpHost = oriExternalHttpHost }()
	externalHttpHost = mockHttpServer.URL()

	assert := require.New(t)
	s := server{}
	req := &sample_service.GetUserByIdReq{
		Uid: int32(externalUid),
	}
	_, err := s.GetUserById(context.Background(), req)
	assert.Error(err)
	assert.Contains(err.Error(), context.DeadlineExceeded.Error())
}

func TestCreateUserFromExternal_WithSuccess(t *testing.T) {
	externalUid := 2
	externalUserInfoPath := fmt.Sprintf("/users/%v", externalUid)
	externalUserInfo := map[string]interface{}{
		"name":   "external user2",
		"gender": 2,
		"age":    20,
	}

	defer func() {
		gock.OffAll()
		gock.DefaultMatcher = gock.NewMatcher()
	}()
	gock.New(externalHttpHost).
		Get(externalUserInfoPath).
		Response.
		Status(http.StatusOK).
		JSON(externalUserInfo)

	newUid := 100
	ctrl := gomock.NewController(t)
	mockGreeterClient := NewMockGreeterClient(ctrl).(*greeter_service.MockGreeterClient)
	mockGreeterClient.EXPECT().
		StoreUserInDb(
			gomock.Any(),
			gomock.Eq(&greeter_service.StoreUserInDbRequest{
				Username:    "external user2",
				Gender:      2,
				Age:         20,
				ExternalUid: int32(externalUid),
			})).
		Return(&greeter_service.StoreUserInDbReply{
			Uid: int64(newUid),
		}, nil)

	assert := require.New(t)
	s := server{}
	req := &sample_service.CreateUserFromExternalReq{
		ExternalUid: int32(externalUid),
	}
	resp, err := s.CreateUserFromExternal(context.Background(), req)
	assert.NoError(err)
	assert.EqualValues(newUid, resp.Uid)
}
