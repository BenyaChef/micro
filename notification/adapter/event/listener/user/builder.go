package usereventlistener

import (
	"errors"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	usecaseinterface "github.com/BenyaChef/micro/notification/boundary/usecase"
)

type Builder struct {
	userEventUseCase usecaseinterface.UserEventUseCase
	logPublisher     loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) UserEventUseCase(userEventUseCase usecaseinterface.UserEventUseCase) *Builder {
	b.userEventUseCase = userEventUseCase

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*UserEventListener, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &UserEventListener{
		userEventUseCase: b.userEventUseCase,
		logPublisher:     b.logPublisher,
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.userEventUseCase == nil {
		errs = append(errs, ErrUserEventUseCaseIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
