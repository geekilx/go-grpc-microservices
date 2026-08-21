package main

import (
	"context"
	"log"

	greetv1 "github.com/geekilx/grpc-course/proto/greet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := greetv1.NewGreetServiceClient(conn)

	response, err := client.Greet(context.Background(), &greetv1.GreetRequest{FirstName: "ilia"})
	if err != nil {
		log.Fatalf("error while calling Greet: %v", err)
	}
	log.Printf("Greet response: %s", response.GetResult())

	log.Println("Calling GreetManyTimes")

	DoGreetManyTimes(client)

	log.Println("Calling LongGreet")

	DoLongGreet(client)

	log.Println("Calling GreetEveryone")
	DoGreetEveryone(client)

}
