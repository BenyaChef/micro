package grpcservermiddleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	infraotel "github.com/BenyaChef/micro/infrastructure/otel"
)

func Metrics(monitor *infraotel.GRPCMonitor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		started := time.Now()

		response, err := handler(ctx, request)

		monitor.RecordRequest(ctx, info.FullMethod, status.Code(err).String(), time.Since(started).Seconds())

		return response, err
	}
}
