package main

import (
	"context"
	"fmt"
	"io"
	"log"

	greetv1 "github.com/geekilx/grpc-course/proto/greet/v1"
)

func DoGreetManyTimes(client greetv1.GreetServiceClient) {

	stream, err := client.GreetManyTimes(context.Background(), &greetv1.GreetRequest{FirstName: "ilia"})
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving stream: %v", err)
		}
		log.Printf("GreetManyTimes response: %s", response.GetResult())
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("error while closing stream: %v", err)
	}

}

func DoLongGreet(client greetv1.GreetServiceClient) {

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for i := 0; i < 10; i++ {
		if err := stream.Send(&greetv1.GreetRequest{FirstName: fmt.Sprintf("ilia %d", i)}); err != nil {
			log.Fatalf("failed to send stream: %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while closing stream: %v", err)
	}
	log.Printf("LongGreet response: %s", resp.GetResult())

}
