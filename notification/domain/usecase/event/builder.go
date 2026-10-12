package eventusecase

import (
	"errors"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	eventusecasehandler "github.com/BenyaChef/micro/notification/domain/usecase/event/handler"
)

type Builder struct {
	handlerProvider *eventusecasehandler.Provider
	logPublisher    loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) HandlerProvider(handlerProvider *eventusecasehandler.Provider) *Builder {
	b.handlerProvider = handlerProvider

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*UserEventUseCase, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &UserEventUseCase{
		handlerProvider: b.handlerProvider,
		logPublisher:    b.logPublisher,
		errorProcessor:  newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.handlerProvider == nil {
		errs = append(errs, ErrHandlerProviderIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
