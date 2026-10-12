package usereventlistenermock

import (
	"context"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	usecaseinterface "github.com/BenyaChef/micro/notification/boundary/usecase"
)

var _ usecaseinterface.UserEventUseCase = (*UserEventUseCaseMock)(nil)

type UserEventUseCaseMock struct {
	HandleFunc func(ctx context.Context, event *commonautheventmodel.UserEvent) error

	Handled []*commonautheventmodel.UserEvent
}

func NewUserEventUseCaseMock() *UserEventUseCaseMock {
	return &UserEventUseCaseMock{
		HandleFunc: func(context.Context, *commonautheventmodel.UserEvent) error { return nil },
	}
}

func (m *UserEventUseCaseMock) Handle(ctx context.Context, event *commonautheventmodel.UserEvent) error {
	m.Handled = append(m.Handled, event)

	return m.HandleFunc(ctx, event)
}
