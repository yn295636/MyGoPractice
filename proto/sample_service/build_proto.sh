#!/usr/bin/env bash
protoc *.proto --go_out=plugins=grpc:./
mockgen -source sample_service.pb.go -package sample_service -destination mock_sample_service.pb.go
