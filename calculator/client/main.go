package main

import (
	"context"
	"io"
	"log"
	"math/rand/v2"

	calculatorv1 "github.com/geekilx/grpc-course/proto/calculator/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	client := calculatorv1.NewCalculatorServiceClient(conn)

	stream, err := client.Prime(context.Background(), &calculatorv1.PrimeRequest{A: 120})
	if err != nil {
		log.Fatalf("failed to call Prime: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive response: %v", err)
		}
		log.Printf("Response received: %v", resp.GetResult())
	}

	sumResp, err := client.Sum(context.Background(), &calculatorv1.SumRequest{A: 3, B: 10})
	if err != nil {
		log.Fatalf("failed to call sum: %v", err)
	}
	log.Printf("Sum response: %v", sumResp.GetResult())

	avgStream, err := client.Avg(context.Background())
	if err != nil {
		log.Fatalf("failed to call Avg: %v", err)
	}
	for i := 1; i < 5; i++ {
		if err := avgStream.Send(&calculatorv1.AvgRequest{A: float32(i)}); err != nil {
			log.Fatalf("failed to send avg request: %v", err)
		}
	}

	avgResp, err := avgStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive avg response: %v", err)
	}

	log.Printf("Avg response: %v", avgResp.GetResult())

	//////////// bidirectional stream ////////////
	maxStream, err := client.Max(context.Background())
	if err != nil {
		log.Fatalf("failed while calling max: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for i := 1; i < 5; i++ {
			if err := maxStream.Send(&calculatorv1.MaxRequest{Num: rand.Int32()}); err != nil {
				log.Fatalf("error while sending numbers to the stream: %v", err)
			}
		}
		maxStream.CloseSend()
	}()

	go func() {
		for {
			res, err := maxStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to receive response: %v", err)
			}
			log.Printf("Response received: %v", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc

}
