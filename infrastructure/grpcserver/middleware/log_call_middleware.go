package grpcservermiddleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

func LogCall(logPublisher loggerinterface.LogPublisher) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		started := time.Now()

		response, err := handler(ctx, request)

		logPublisher.LogInfo(ctx, "grpc call handled",
			"method", info.FullMethod,
			"code", status.Code(err).String(),
			"took", time.Since(started),
		)

		return response, err
	}
}
