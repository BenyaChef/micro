package logger

import (
	"errors"
	"log/slog"
	"os"

	loggermodel "github.com/BenyaChef/micro/infrastructure/logger/model"
)

type Builder struct {
	serviceName string
	level       string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) ServiceName(serviceName string) *Builder {
	b.serviceName = serviceName

	return b
}

func (b *Builder) Level(level string) *Builder {
	b.level = level

	return b
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.serviceName == "" {
		errs = append(errs, ErrServiceNameIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (b *Builder) Build() (*LogPublisher, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggermodel.ParseLevel(b.level),
	})

	return &LogPublisher{
		log: slog.New(handler).With(slog.String("service", b.serviceName)),
	}, nil
}
