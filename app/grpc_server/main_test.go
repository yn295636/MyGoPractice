package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	pb "github.com/yn295636/MyGoPractice/proto"
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
}
