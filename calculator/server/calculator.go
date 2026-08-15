package main

import (
	calculatorv1 "github.com/geekilx/grpc-course/proto/calculator/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalculatorService struct {
	calculatorv1.UnimplementedCalculatorServiceServer
}

func NewCalculatorSerivce() *CalculatorService {
	return &CalculatorService{}
}

func (s *CalculatorService) Prime(req *calculatorv1.PrimeRequest, stream calculatorv1.CalculatorService_PrimeServer) error {

	if req.GetA() == 0 {
		return status.Error(codes.InvalidArgument, "base number is required")
	}

	k := 2
	n := int(req.GetA())
	for n > 1 {
		if n%k == 0 {
			stream.Send(&calculatorv1.PrimeResponse{Result: int32(k)})
			n = n / k
		} else {
			k++
		}
	}
	return nil

}
