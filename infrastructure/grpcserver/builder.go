package grpcserver

import (
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"

	grpcservermiddleware "github.com/BenyaChef/micro/infrastructure/grpcserver/middleware"
	grpcservermodel "github.com/BenyaChef/micro/infrastructure/grpcserver/model"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Builder struct {
	logPublisher    loggerinterface.LogPublisher
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
		grpc.UnaryInterceptor(grpcservermiddleware.LogCall(b.logPublisher)),
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
