package userusecasemock

import (
	"context"

	eventinterface "github.com/BenyaChef/micro/auth/boundary/event"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
)

var _ eventinterface.UserEventPublisher = (*UserEventPublisherMock)(nil)

type UserEventPublisherMock struct {
	RegisteredFunc func(ctx context.Context, userID commonuserentity.UserID, email emailprimitive.Email) error

	Calls      int
	LastUserID commonuserentity.UserID
	LastEmail  emailprimitive.Email
}

func NewUserEventPublisherMock() *UserEventPublisherMock {
	return &UserEventPublisherMock{
		RegisteredFunc: func(context.Context, commonuserentity.UserID, emailprimitive.Email) error {
			return nil
		},
	}
}

func (m *UserEventPublisherMock) Registered(
	ctx context.Context,
	userID commonuserentity.UserID,
	email emailprimitive.Email,
) error {
	m.Calls++
	m.LastUserID = userID
	m.LastEmail = email

	return m.RegisteredFunc(ctx, userID, email)
}
