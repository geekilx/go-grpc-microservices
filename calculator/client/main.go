package main

import (
	"context"
	"io"
	"log"

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

}
