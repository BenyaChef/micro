package authrestcontroller

import (
	"errors"

	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
)

type Builder struct {
	baseController *restservercontroller.BaseController
	userUseCase    usecaseinterface.UserUseCase
	tokenUseCase   usecaseinterface.TokenUseCase
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BaseController(baseController *restservercontroller.BaseController) *Builder {
	b.baseController = baseController

	return b
}

func (b *Builder) UserUseCase(userUseCase usecaseinterface.UserUseCase) *Builder {
	b.userUseCase = userUseCase

	return b
}

func (b *Builder) TokenUseCase(tokenUseCase usecaseinterface.TokenUseCase) *Builder {
	b.tokenUseCase = tokenUseCase

	return b
}

func (b *Builder) Build() (*AuthController, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &AuthController{
		BaseController: b.baseController,
		userUseCase:    b.userUseCase,
		tokenUseCase:   b.tokenUseCase,
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.baseController == nil {
		errs = append(errs, ErrBaseControllerIsRequired)
	}

	if b.userUseCase == nil {
		errs = append(errs, ErrUserUseCaseIsRequired)
	}

	if b.tokenUseCase == nil {
		errs = append(errs, ErrTokenUseCaseIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
