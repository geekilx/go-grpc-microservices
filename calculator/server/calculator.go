package main

import (
	"context"
	"io"
	"log"

	calculatorv1 "github.com/geekilx/go-grpc-microservices/proto/calculator/v1"
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

func (s *CalculatorService) Sum(ctx context.Context, req *calculatorv1.SumRequest) (*calculatorv1.SumResponse, error) {

	if req.GetA() == 0 || req.GetB() == 0 {
		return nil, status.Error(codes.InvalidArgument, "base numbers are required")
	}

	return &calculatorv1.SumResponse{Result: req.GetA() + req.GetB()}, nil

}

func (s *CalculatorService) Avg(stream calculatorv1.CalculatorService_AvgServer) error {

	var sum float32
	var count float32
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sum += resp.GetA()
		count++
	}

	return stream.SendAndClose(&calculatorv1.AvgResponse{Result: sum / count})

}

func (s *CalculatorService) Max(stream calculatorv1.CalculatorService_MaxServer) error {

	var num int32
	var allNums []int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "failed to receive max request")
		}

		log.Printf("received number from client: %d", req.GetNum())

		if num < req.GetNum() {
			num = req.GetNum()
			if err := stream.Send(&calculatorv1.MaxResponse{Result: num}); err != nil {
				return status.Error(codes.Internal, "failed to send max response")
			}
			allNums = append(allNums, num)
		}

	}

	return nil

}
