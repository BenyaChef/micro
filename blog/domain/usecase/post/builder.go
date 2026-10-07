package postusecase

import (
	"errors"

	repositoryinterface "github.com/BenyaChef/micro/blog/boundary/repository"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Builder struct {
	postRepository repositoryinterface.PostRepository
	logPublisher   loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) PostRepository(postRepository repositoryinterface.PostRepository) *Builder {
	b.postRepository = postRepository

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*PostUseCase, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &PostUseCase{
		postRepository: b.postRepository,
		logPublisher:   b.logPublisher,
		errorProcessor: newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.postRepository == nil {
		errs = append(errs, ErrPostRepositoryIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
