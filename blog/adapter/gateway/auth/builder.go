package authgateway

import (
	"errors"

	commonauthgatewayinterface "github.com/BenyaChef/micro/common/auth/boundary/gateway"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Builder struct {
	authGateway  commonauthgatewayinterface.AuthGateway
	logPublisher loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) AuthGateway(authGateway commonauthgatewayinterface.AuthGateway) *Builder {
	b.authGateway = authGateway

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*AuthGateway, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &AuthGateway{
		authGateway:  b.authGateway,
		logPublisher: b.logPublisher,
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.authGateway == nil {
		errs = append(errs, ErrAuthGatewayIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
