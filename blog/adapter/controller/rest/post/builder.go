package postrestcontroller

import (
	"errors"

	usecaseinterface "github.com/BenyaChef/micro/blog/boundary/usecase"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
)

type Builder struct {
	baseController *restservercontroller.BaseController
	postUseCase    usecaseinterface.PostUseCase
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BaseController(baseController *restservercontroller.BaseController) *Builder {
	b.baseController = baseController

	return b
}

func (b *Builder) PostUseCase(postUseCase usecaseinterface.PostUseCase) *Builder {
	b.postUseCase = postUseCase

	return b
}

func (b *Builder) Build() (*PostController, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &PostController{
		BaseController: b.baseController,
		postUseCase:    b.postUseCase,
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.baseController == nil {
		errs = append(errs, ErrBaseControllerIsRequired)
	}

	if b.postUseCase == nil {
		errs = append(errs, ErrPostUseCaseIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
