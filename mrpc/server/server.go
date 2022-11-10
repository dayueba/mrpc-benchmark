package main

import (
	"time"

	"github.com/dayueba/mrpc-benchmark/testdata"

	"github.com/dayueba/mrpc"
)


func main() {
	opts := []mrpc.ServerOption{
		mrpc.WithAddress("127.0.0.1:8000"),
		mrpc.WithTimeout(time.Millisecond * 20000),
	}
	s := mrpc.NewServer(opts ...)


	if err := s.RegisterService("helloworld.Greeter", new(testdata.Service)); err != nil {
		panic(err)
	}

	s.Serve()
}