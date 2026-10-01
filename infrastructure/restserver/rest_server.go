package restserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	restserverservice "github.com/BenyaChef/micro/infrastructure/restserver/service"
)

type Server struct {
	logPublisher    loggerinterface.LogPublisher
	responseService *restserverservice.ResponseService
	router          chi.Router
	server          *http.Server
	shutdownTimeout time.Duration
}

func (s *Server) Router() chi.Router {
	return s.router
}

func (s *Server) ResponseService() *restserverservice.ResponseService {
	return s.responseService
}

func (s *Server) Addr() string {
	return s.server.Addr
}

func (s *Server) Run(ctx context.Context) error {
	listenFailed := make(chan error, 1)

	go func() {
		s.logPublisher.LogInfo(ctx, "rest server listening", "addr", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			listenFailed <- ErrListenFailed(s.server.Addr, err)
		}

		close(listenFailed)
	}()

	select {
	case err := <-listenFailed:
		return err
	case <-ctx.Done():
		s.logPublisher.LogInfo(ctx, "rest server shutting down")
	}
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return ErrShutdownFailed(s.server.Addr, err)
	}

	s.logPublisher.LogInfo(shutdownCtx, "rest server stopped")

	return nil
}
