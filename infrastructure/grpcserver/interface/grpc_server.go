package grpcserverinterface

import (
	"context"

	"google.golang.org/grpc"
)

type GRPCServer interface {
	Server() *grpc.Server
	Addr() string
	Run(ctx context.Context) error
}
