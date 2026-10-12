package eventusecasemock

import (
	"context"

	repositoryinterface "github.com/BenyaChef/micro/notification/boundary/repository"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
)

var _ repositoryinterface.NotificationRepository = (*NotificationRepositoryMock)(nil)

type NotificationRepositoryMock struct {
	InsertFunc       func(ctx context.Context, notification *notificationentity.Notification) (bool, error)
	GetByEventIDFunc func(ctx context.Context, eventID notificationentity.EventID) (*notificationentity.Notification, error)

	Inserted []*notificationentity.Notification
}

func NewNotificationRepositoryMock() *NotificationRepositoryMock {
	return &NotificationRepositoryMock{
		InsertFunc: func(context.Context, *notificationentity.Notification) (bool, error) { return true, nil },
		GetByEventIDFunc: func(context.Context, notificationentity.EventID) (*notificationentity.Notification, error) {
			return nil, nil
		},
	}
}

func (m *NotificationRepositoryMock) Insert(
	ctx context.Context,
	notification *notificationentity.Notification,
) (bool, error) {
	m.Inserted = append(m.Inserted, notification)

	return m.InsertFunc(ctx, notification)
}

func (m *NotificationRepositoryMock) GetByEventID(
	ctx context.Context,
	eventID notificationentity.EventID,
) (*notificationentity.Notification, error) {
	return m.GetByEventIDFunc(ctx, eventID)
}
