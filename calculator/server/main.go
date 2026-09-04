package main

import (
	"log"
	"net"

	calculatorv1 "github.com/geekilx/go-grpc-microservices/proto/calculator/v1"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("server listening on :50051")

	s := grpc.NewServer()
	grpcService := NewCalculatorSerivce()
	calculatorv1.RegisterCalculatorServiceServer(s, grpcService)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
