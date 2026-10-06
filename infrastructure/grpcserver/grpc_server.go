package grpcserver

import (
	"context"
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"

	grpcserverinterface "github.com/BenyaChef/micro/infrastructure/grpcserver/interface"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const network = "tcp"

var _ grpcserverinterface.GRPCServer = (*Server)(nil)

type Server struct {
	logPublisher    loggerinterface.LogPublisher
	server          *grpc.Server
	addr            string
	shutdownTimeout time.Duration
}

func (s *Server) Server() *grpc.Server {
	return s.server
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Run(ctx context.Context) error {
	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(ctx, network, s.addr)
	if err != nil {
		return ErrListenFailed(s.addr, err)
	}

	serveFailed := make(chan error, 1)

	go func() {
		s.logPublisher.LogInfo(ctx, "grpc server listening", "addr", s.addr)

		if err := s.server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serveFailed <- ErrServeFailed(s.addr, err)
		}
	}()

	select {
	case err := <-serveFailed:
		return err
	case <-ctx.Done():
		s.logPublisher.LogInfo(ctx, "grpc server shutting down")
	}

	return s.stop(ctx)
}

func (s *Server) stop(ctx context.Context) error {
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.shutdownTimeout)
	defer cancel()

	select {
	case <-stopped:
		s.logPublisher.LogInfo(shutdownCtx, "grpc server stopped")

		return nil
	case <-shutdownCtx.Done():
		s.server.Stop()

		return ErrShutdownTimeout(s.addr, s.shutdownTimeout)
	}
}
