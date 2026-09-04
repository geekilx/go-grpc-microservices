package main

import (
	"log"
	"net"

	greetv1 "github.com/geekilx/go-grpc-microservices/proto/greet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("listening on :50051")

	s := grpc.NewServer()
	grpcServer := NewGreetService()

	greetv1.RegisterGreetServiceServer(s, grpcServer)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
