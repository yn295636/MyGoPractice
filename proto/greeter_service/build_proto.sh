#!/usr/bin/env bash
protoc *.proto --go_out=plugins=grpc:./
mockgen -source greeter_service.pb.go -package greeter_service -destination mock_greeter_service.pb.go
