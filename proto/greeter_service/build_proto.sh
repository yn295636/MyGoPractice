#!/usr/bin/env bash
protoc *.proto -I=./ --go_out=./ --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:./ --go-grpc_opt=paths=source_relative
mockgen -source greeter_service.pb.go -package greeter_service -destination mock_greeter_service.pb.go
mockgen -source greeter_service_grpc.pb.go -package greeter_service -destination mock_greeter_service_grpc.pb.go
