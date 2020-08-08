package main

import (
	"context"
	"github.com/yn295636/MyGoPractice/proto/sample_service"
)

type server struct{}

func (s *server) Multiply(ctx context.Context, in *sample_service.MultiplyReq) (*sample_service.MultiplyResp, error) {
	return &sample_service.MultiplyResp{
		Result: int64(in.A) * int64(in.B),
	}, nil
}
