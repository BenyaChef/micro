package userrepository

import (
	"errors"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	postgresinterface "github.com/BenyaChef/micro/infrastructure/postgres/interface"
)

type Builder struct {
	executor     postgresinterface.Executor
	logPublisher loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Executor(executor postgresinterface.Executor) *Builder {
	b.executor = executor

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*UserRepository, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &UserRepository{
		executor:       b.executor,
		logPublisher:   b.logPublisher,
		errorProcessor: newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.executor == nil {
		errs = append(errs, ErrExecutorIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
