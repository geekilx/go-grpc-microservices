package main

import (
	"context"
	"fmt"

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
