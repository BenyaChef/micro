package tokenusecase

import (
	"errors"

	serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Builder struct {
	tokenService serviceinterface.TokenService
	logPublisher loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) TokenService(tokenService serviceinterface.TokenService) *Builder {
	b.tokenService = tokenService

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*TokenUseCase, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &TokenUseCase{
		tokenService:   b.tokenService,
		logPublisher:   b.logPublisher,
		errorProcessor: newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.tokenService == nil {
		errs = append(errs, ErrTokenServiceIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
