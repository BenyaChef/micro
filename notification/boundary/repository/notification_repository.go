package repositoryinterface

import (
	"context"

	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
)

type NotificationRepository interface {
	Insert(ctx context.Context, notification *notificationentity.Notification) (inserted bool, err error)
	GetByEventID(ctx context.Context, eventID notificationentity.EventID) (*notificationentity.Notification, error)
}
