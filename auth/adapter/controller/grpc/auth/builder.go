package authgrpccontroller

import (
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
)

type Builder struct {
	tokenUseCase usecaseinterface.TokenUseCase
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) TokenUseCase(tokenUseCase usecaseinterface.TokenUseCase) *Builder {
	b.tokenUseCase = tokenUseCase

	return b
}

func (b *Builder) Build() (*AuthController, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &AuthController{tokenUseCase: b.tokenUseCase}, nil
}

func (b *Builder) checkRequiredFields() error {
	if b.tokenUseCase == nil {
		return ErrTokenUseCaseIsRequired
	}

	return nil
}
