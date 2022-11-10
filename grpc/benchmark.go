package main

import (
	"context"
	"flag"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/dayueba/mrpc-benchmark/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"github.com/dayueba/mrpc/log"
)

var concurrency = flag.Int64("concurrency", 500, "concurrency")
var target = flag.String("target", "127.0.0.1:8000", "target")
var total = flag.Int64("total", 1000000, "total requests")

const addr = ":50051"

func main() {
	flag.Parse()
	request(*total, *concurrency, *target)
}

type Counter struct {
	Succ        int64 // 成功量
	Fail        int64 // 失败量
	Total       int64 // 总量
	Concurrency int64 // 并发量
	Cost        int64 // 总耗时 ms
}

func request(totalReqs int64, concurrency int64, target string) {
	perClientReqs := totalReqs / concurrency
	counter := &Counter{
		Total:       perClientReqs * concurrency,
		Concurrency: concurrency,
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)
	
	var wg sync.WaitGroup
	wg.Add(int(concurrency))

	startTime := time.Now().UnixNano()

	for i := int64(0); i < counter.Concurrency; i++ {
		go func(i int64) {
			for j := int64(0); j < perClientReqs; j++ {
				_, err := c.Hello(context.Background(), &pb.Request{FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"hello", "world"},
				}})
				if err != nil {
					log.Info(err)
					atomic.AddInt64(&counter.Fail, 1)
				} else {
					atomic.AddInt64(&counter.Succ, 1)
				}			
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	counter.Cost = (time.Now().UnixNano() - startTime) / 1000000

	log.Infof("took %d ms for %d requests", counter.Cost, counter.Total)
	log.Infof("sent     requests      : %d\n", counter.Total)
	log.Infof("received requests      : %d\n", atomic.LoadInt64(&counter.Succ)+atomic.LoadInt64(&counter.Fail))
	log.Infof("received requests succ : %d\n", atomic.LoadInt64(&counter.Succ))
	log.Infof("received requests fail : %d\n", atomic.LoadInt64(&counter.Fail))
	log.Infof("throughput  (TPS)      : %d\n", totalReqs*1000/counter.Cost)
}
