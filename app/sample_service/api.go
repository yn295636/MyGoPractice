package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yn295636/MyGoPractice/proto/greeter_service"
	"github.com/yn295636/MyGoPractice/proto/sample_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"time"
)

var (
	externalHttpHost = "http://192.168.56.1:5000"
)

type ExternalErrorResp struct {
	ErrorCode int32  `json:"error"`
	Message   string `json:"message"`
}

type ExternalUserInfo struct {
	Name   string `json:"name"`
	Gender uint32 `json:"gender"`
	Age    uint32 `json:"age"`
}

type server struct{}

func (s *server) Multiply(ctx context.Context, in *sample_service.MultiplyReq) (*sample_service.MultiplyResp, error) {
	return &sample_service.MultiplyResp{
		Result: int64(in.A) * int64(in.B),
	}, nil
}

func (s *server) GetUserById(ctx context.Context,
	in *sample_service.GetUserByIdReq) (*sample_service.GetUserByIdResp, error) {
	userInfoApiPath := fmt.Sprintf("/users/%v", in.Uid)
	log.Printf("Will get external user info thru url %v%v", externalHttpHost, userInfoApiPath)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%v%v", externalHttpHost, userInfoApiPath), nil)
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error happens for the external http request: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var externalUserInfo ExternalUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&externalUserInfo); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Parsing external user info from response body got failed! %v", err.Error()))
		}
		log.Printf("external user info by uid %v: %v", in.Uid, externalUserInfo)
		return &sample_service.GetUserByIdResp{
			Name:   externalUserInfo.Name,
			Gender: sample_service.GetUserByIdResp_Gender(externalUserInfo.Gender),
			Age:    externalUserInfo.Age,
		}, nil
	case http.StatusBadRequest:
		var externalError ExternalErrorResp
		if err := json.NewDecoder(resp.Body).Decode(&externalError); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Parsing external error from response body got failed! %v", err.Error()))
		}
		log.Printf("http status: %v, external error: %v", resp.StatusCode, externalError)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("external message: %v", externalError.Message))
	default:
		var externalError ExternalErrorResp
		if err := json.NewDecoder(resp.Body).Decode(&externalError); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Parsing external error from response body got failed! %v", err.Error()))
		}
		log.Printf("http status: %v, external error: %v", resp.StatusCode, externalError)
		return nil, status.Error(codes.Unavailable, fmt.Sprintf("external error code: %v, external message: %v", externalError.ErrorCode, externalError.Message))
	}
}

func (s *server) CreateUserFromExternal(ctx context.Context, in *sample_service.CreateUserFromExternalReq) (*sample_service.CreateUserFromExternalResp, error) {
	externalUserInfo, err := s.GetUserById(ctx, &sample_service.GetUserByIdReq{Uid: in.ExternalUid})
	if err != nil {
		log.Printf("Get external user info failed: %v", err)
		return nil, err
	}

	// Get client for greeter_service
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("greeter_service cannot be connected: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	greeterClient := greeter_service.NewGreeterClient(conn)
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatalf("release connection error: %v", err)
		}
	}()

	//greeterClient, err, releaseConn := GetGreeterClient()
	//if err != nil {
	//	log.Printf("Get greeter_service client got error: %v", err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//defer releaseConn()

	resp, err := greeterClient.StoreUserInDb(ctx, &greeter_service.StoreUserInDbRequest{
		Username:    externalUserInfo.Name,
		Gender:      int32(externalUserInfo.Gender),
		Age:         int32(externalUserInfo.Age),
		ExternalUid: in.ExternalUid,
	})
	if err != nil {
		log.Printf("Create user got error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("Create user successfully and returned uid: %v", resp.Uid)
	return &sample_service.CreateUserFromExternalResp{
		Uid: resp.Uid,
	}, nil
}
