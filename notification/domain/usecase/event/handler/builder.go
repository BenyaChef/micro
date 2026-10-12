package eventusecasehandler

import (
	"errors"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	repositoryinterface "github.com/BenyaChef/micro/notification/boundary/repository"
)

type Builder struct {
	notificationRepository repositoryinterface.NotificationRepository
	logPublisher           loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) NotificationRepository(repository repositoryinterface.NotificationRepository) *Builder {
	b.notificationRepository = repository

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) BuildRegisteredHandler() (*RegisteredHandler, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &RegisteredHandler{
		notificationRepository: b.notificationRepository,
		logPublisher:           b.logPublisher,
		errorProcessor:         newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.notificationRepository == nil {
		errs = append(errs, ErrNotificationRepositoryIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
