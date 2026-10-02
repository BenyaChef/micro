package authrestcontrollermock

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

var _ usecaseinterface.UserUseCase = (*UserUseCaseMock)(nil)

type UserUseCaseMock struct {
	RegisterFunc func(ctx context.Context, data *boundarymodel.RegisterData) (*userentity.User, error)
	LoginFunc    func(ctx context.Context, data *boundarymodel.LoginData) (*boundarymodel.LoginResult, error)
	GetByIDFunc  func(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error)
}

func NewUserUseCaseMock() *UserUseCaseMock {
	return &UserUseCaseMock{
		RegisterFunc: func(context.Context, *boundarymodel.RegisterData) (*userentity.User, error) {
			return nil, nil
		},
		LoginFunc: func(context.Context, *boundarymodel.LoginData) (*boundarymodel.LoginResult, error) {
			return nil, nil
		},
		GetByIDFunc: func(context.Context, commonuserentity.UserID) (*userentity.User, error) {
			return nil, nil
		},
	}
}

func (m *UserUseCaseMock) Register(ctx context.Context, data *boundarymodel.RegisterData) (*userentity.User, error) {
	return m.RegisterFunc(ctx, data)
}

func (m *UserUseCaseMock) Login(ctx context.Context, data *boundarymodel.LoginData) (*boundarymodel.LoginResult, error) {
	return m.LoginFunc(ctx, data)
}

func (m *UserUseCaseMock) GetByID(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error) {
	return m.GetByIDFunc(ctx, userID)
}
