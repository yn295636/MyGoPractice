#!/usr/bin/env bash
protoc *.proto -I=./ --go_out=./ --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:./ --go-grpc_opt=paths=source_relative
mockgen -source sample_service.pb.go -package sample_service -destination mock_sample_service.pb.go
mockgen -source sample_service_grpc.pb.go -package sample_service -destination mock_sample_service_grpc.pb.go
