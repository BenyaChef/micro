package commonauthgateway

import (
	"errors"

	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Builder struct {
	client       authproto.AuthServiceClient
	logPublisher loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Client(client authproto.AuthServiceClient) *Builder {
	b.client = client

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*AuthGRPCGateway, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &AuthGRPCGateway{
		client:         b.client,
		logPublisher:   b.logPublisher,
		errorProcessor: newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.client == nil {
		errs = append(errs, ErrClientIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
