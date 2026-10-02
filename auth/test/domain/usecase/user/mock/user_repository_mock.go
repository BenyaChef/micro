package userusecasemock

import (
	"context"

	repositoryinterface "github.com/BenyaChef/micro/auth/boundary/repository"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
)

var _ repositoryinterface.UserRepository = (*UserRepositoryMock)(nil)

type UserRepositoryMock struct {
	InsertFunc     func(ctx context.Context, user *userentity.User) error
	GetByIDFunc    func(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error)
	GetByEmailFunc func(ctx context.Context, email emailprimitive.Email) (*userentity.User, error)

	Inserted *userentity.User
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{
		InsertFunc:     func(context.Context, *userentity.User) error { return nil },
		GetByIDFunc:    func(context.Context, commonuserentity.UserID) (*userentity.User, error) { return nil, nil },
		GetByEmailFunc: func(context.Context, emailprimitive.Email) (*userentity.User, error) { return nil, nil },
	}
}

func (m *UserRepositoryMock) Insert(ctx context.Context, user *userentity.User) error {
	m.Inserted = user

	return m.InsertFunc(ctx, user)
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error) {
	return m.GetByIDFunc(ctx, userID)
}

func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email emailprimitive.Email) (*userentity.User, error) {
	return m.GetByEmailFunc(ctx, email)
}
