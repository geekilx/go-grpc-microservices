package main

import (
	"context"
	"time"

	"google.golang.org/grpc"
)



func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req,replay any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts...grpc.CallOption) error {

		timeoutCtx,cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(timeoutCtx, method, req, replay, cc, opts...)

	}
}

func StreamTimeoutInterceptor(timeout time.Duration) grpc.StreamClientInterceptor {
return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)

	s, err := streamer(timeoutCtx, desc, cc, method, opts...)
	if err != nil {
		cancel()
		return nil, err
	}


	return &cancelableStream{ClientStream: s, cancel: cancel}, nil
}
}

type cancelableStream struct {
	grpc.ClientStream
	cancel context.CancelFunc
}

func (cs *cancelableStream) RecvMsg(m any) error{

	err := cs.ClientStream.RecvMsg(m)
	if err != nil {
		cs.cancel()
	}

	return err
}

