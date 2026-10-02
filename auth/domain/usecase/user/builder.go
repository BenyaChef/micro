package userusecase

import (
	"errors"

	repositoryinterface "github.com/BenyaChef/micro/auth/boundary/repository"
	serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Builder struct {
	userRepository repositoryinterface.UserRepository
	tokenUseCase   usecaseinterface.TokenUseCase
	passwordHasher serviceinterface.PasswordHasher
	logPublisher   loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) UserRepository(userRepository repositoryinterface.UserRepository) *Builder {
	b.userRepository = userRepository

	return b
}

func (b *Builder) TokenUseCase(tokenUseCase usecaseinterface.TokenUseCase) *Builder {
	b.tokenUseCase = tokenUseCase

	return b
}

func (b *Builder) PasswordHasher(passwordHasher serviceinterface.PasswordHasher) *Builder {
	b.passwordHasher = passwordHasher

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build() (*UserUseCase, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &UserUseCase{
		userRepository: b.userRepository,
		tokenUseCase:   b.tokenUseCase,
		passwordHasher: b.passwordHasher,
		logPublisher:   b.logPublisher,
		errorProcessor: newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.userRepository == nil {
		errs = append(errs, ErrUserRepositoryIsRequired)
	}

	if b.tokenUseCase == nil {
		errs = append(errs, ErrTokenUseCaseIsRequired)
	}

	if b.passwordHasher == nil {
		errs = append(errs, ErrPasswordHasherIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
