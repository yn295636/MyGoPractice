#!/usr/bin/env bash
protoc *.proto --go_out=plugins=grpc:./
mockgen -source helloworld.pb.go -package proto -destination mockHelloworld.pb.go
