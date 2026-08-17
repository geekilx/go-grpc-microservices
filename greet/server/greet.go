package main

import (
	"context"
	"fmt"
	"io"
	"log"

	greetv1 "github.com/geekilx/grpc-course/proto/greet/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreetService struct {
	greetv1.UnimplementedGreetServiceServer
}

func NewGreetService() *GreetService {
	return &GreetService{}
}

func (s *GreetService) Greet(ctx context.Context, req *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {

	if req.GetFirstName() == "" {
		return nil, status.Error(codes.InvalidArgument, "First name is required")
	}

	return &greetv1.GreetResponse{
		Result: fmt.Sprintf("Hello from %v", req.GetFirstName()),
	}, nil

}

func (s *GreetService) GreetManyTimes(req *greetv1.GreetRequest, stream greetv1.GreetService_GreetManyTimesServer) error {

	log.Printf("GreetManyTimes called with %v", req)

	for i := 0; i < 10; i++ {
		if err := stream.Send(&greetv1.GreetResponse{Result: fmt.Sprintf("Hello from %v, number %d", req.GetFirstName(), i)}); err != nil {
			log.Fatalf("failed to send stream: %v", err)
		}
	}

	return nil

}

func (s *GreetService) LongGreet(stream greetv1.GreetService_LongGreetServer) error {

	log.Printf("LongGreet called")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetv1.GreetResponse{Result: "i got all requests :)"})
		}
		if err != nil {
			return err
		}
		log.Printf("LongGreet response: %s", req.GetFirstName())
	}

}
