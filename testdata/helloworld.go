package testdata

import (
	"context"
	// "errors"
	"sync/atomic"
)

type DB string

type Service struct {
	count int64
}

type HelloRequest struct {
	Msg string
}

type HelloReply struct {
	Msg string
}

type AddRequest struct {
	A int32 `msgpack:"a"`
	B int32 `msgpack:"b"`
}

type AddReply struct {
	Result int32 `msgpack:"result"`
}

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	rsp := &HelloReply{
		Msg : "world",
	}

	atomic.AddInt64(&s.count, 1)

	return rsp, nil
}

func (s *Service) Count(ctx context.Context, req *HelloRequest) (*CountResponse, error) {
	atomic.AddInt64(&s.count, 1)
	// fmt.Println(s.count)

	rsp := &CountResponse{
		Count: 100,
	}
	return rsp, nil
}
