package otel

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	network                = "tcp"
	defaultShutdownTimeout = 5 * time.Second
	metricsPath            = "/metrics"
)

type MetricsServer struct {
	server       *http.Server
	handler      http.Handler
	logPublisher loggerinterface.LogPublisher
	addr         string
}

func NewMetricsServer(handler http.Handler, port string, logPublisher loggerinterface.LogPublisher) *MetricsServer {
	mux := http.NewServeMux()
	mux.Handle(metricsPath, handler)

	addr := net.JoinHostPort("", port)

	return &MetricsServer{
		server:       &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: defaultShutdownTimeout},
		handler:      handler,
		logPublisher: logPublisher,
		addr:         addr,
	}
}

func (s *MetricsServer) Handler() http.Handler {
	return s.handler
}

func (s *MetricsServer) Addr() string {
	return s.addr
}

func (s *MetricsServer) Run(ctx context.Context) error {
	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(ctx, network, s.addr)
	if err != nil {
		return ErrListenFailed(s.addr, err)
	}

	serveFailed := make(chan error, 1)

	go func() {
		s.logPublisher.LogInfo(ctx, "metrics server listening", "addr", s.addr, "path", metricsPath)

		if serveErr := s.server.Serve(listener); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			serveFailed <- ErrListenFailed(s.addr, serveErr)
		}
	}()

	select {
	case err := <-serveFailed:
		return err
	case <-ctx.Done():
		s.logPublisher.LogInfo(ctx, "metrics server shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), defaultShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	s.logPublisher.LogInfo(shutdownCtx, "metrics server stopped")

	return nil
}
