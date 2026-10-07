package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func callTimeout(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		request, response any,
		connection *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if _, hasDeadline := ctx.Deadline(); hasDeadline {
			return invoker(ctx, method, request, response, connection, opts...)
		}

		callCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(callCtx, method, request, response, connection, opts...)
	}
}
