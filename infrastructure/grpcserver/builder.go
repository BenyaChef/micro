package grpcserver

import (
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"

	grpcservermiddleware "github.com/BenyaChef/micro/infrastructure/grpcserver/middleware"
	grpcservermodel "github.com/BenyaChef/micro/infrastructure/grpcserver/model"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	infraotel "github.com/BenyaChef/micro/infrastructure/otel"
)

type Builder struct {
	logPublisher    loggerinterface.LogPublisher
	grpcMonitor     *infraotel.GRPCMonitor
	port            string
	shutdownTimeout time.Duration
}

func NewBuilder() *Builder {
	return &Builder{shutdownTimeout: grpcservermodel.DefaultShutdownTimeout}
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) GRPCMonitor(grpcMonitor *infraotel.GRPCMonitor) *Builder {
	b.grpcMonitor = grpcMonitor

	return b
}

func (b *Builder) Port(port string) *Builder {
	b.port = port

	return b
}

func (b *Builder) ShutdownTimeout(timeout time.Duration) *Builder {
	b.shutdownTimeout = timeout

	return b
}

func (b *Builder) Build() (*Server, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcservermiddleware.LogCall(b.logPublisher),
			grpcservermiddleware.Metrics(b.grpcMonitor),
		),
	)

	return &Server{
		logPublisher:    b.logPublisher,
		server:          server,
		addr:            net.JoinHostPort("", b.port),
		shutdownTimeout: b.shutdownTimeout,
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if b.port == "" {
		errs = append(errs, ErrPortIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
