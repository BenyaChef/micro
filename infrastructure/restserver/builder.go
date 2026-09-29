package restserver

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	loggerinterface "micro/infrastructure/logger/interface"
	restserverhandler "micro/infrastructure/restserver/handler"
	restserverinterface "micro/infrastructure/restserver/interface"
	restservermiddleware "micro/infrastructure/restserver/middleware"
	restservermodel "micro/infrastructure/restserver/model"
	restserverservice "micro/infrastructure/restserver/service"
)

type Builder struct {
	logPublisher    loggerinterface.LogPublisher
	errorResolver   restserverinterface.ErrorResolver
	serviceName     string
	port            string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
}

func NewBuilder() *Builder {
	return &Builder{
		readTimeout:     restservermodel.DefaultReadTimeout,
		writeTimeout:    restservermodel.DefaultWriteTimeout,
		idleTimeout:     restservermodel.DefaultIdleTimeout,
		shutdownTimeout: restservermodel.DefaultShutdownTimeout,
	}
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) ErrorResolver(errorResolver restserverinterface.ErrorResolver) *Builder {
	b.errorResolver = errorResolver

	return b
}

func (b *Builder) ServiceName(serviceName string) *Builder {
	b.serviceName = serviceName

	return b
}

func (b *Builder) Port(port string) *Builder {
	b.port = port

	return b
}

func (b *Builder) ReadTimeout(timeout time.Duration) *Builder {
	b.readTimeout = timeout

	return b
}

func (b *Builder) WriteTimeout(timeout time.Duration) *Builder {
	b.writeTimeout = timeout

	return b
}

func (b *Builder) IdleTimeout(timeout time.Duration) *Builder {
	b.idleTimeout = timeout

	return b
}

func (b *Builder) ShutdownTimeout(timeout time.Duration) *Builder {
	b.shutdownTimeout = timeout

	return b
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.logPublisher == nil {
		errs = append(errs, ErrLoggerIsRequired)
	}

	if b.serviceName == "" {
		errs = append(errs, ErrServiceNameIsRequired)
	}

	if b.port == "" {
		errs = append(errs, ErrPortIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (b *Builder) Build() (*Server, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	errorResolver := b.errorResolver
	if errorResolver == nil {
		errorResolver = restserverservice.NewDefaultErrorResolver()
	}

	errorResponseService, err := restserverservice.NewErrorResponseService(errorResolver, b.logPublisher)
	if err != nil {
		return nil, err
	}

	responseService, err := restserverservice.NewResponseService(errorResponseService, b.logPublisher)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(restservermiddleware.RequestID)
	router.Use(restservermiddleware.Recover(b.logPublisher))
	router.Use(restservermiddleware.LogRequest(b.logPublisher))

	restserverhandler.RegisterHealth(router, b.serviceName, b.logPublisher)

	return &Server{
		logPublisher:    b.logPublisher,
		responseService: responseService,
		router:          router,
		shutdownTimeout: b.shutdownTimeout,
		server: &http.Server{
			Addr:         net.JoinHostPort("", b.port),
			Handler:      router,
			ReadTimeout:  b.readTimeout,
			WriteTimeout: b.writeTimeout,
			IdleTimeout:  b.idleTimeout,
		},
	}, nil
}
