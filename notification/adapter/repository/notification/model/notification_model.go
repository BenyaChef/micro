package notificationrepomodel

import (
	"time"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
)

type Notification struct {
	ID      string    `db:"id"`
	EventID string    `db:"event_id"`
	UserID  string    `db:"user_id"`
	Kind    string    `db:"kind"`
	Payload []byte    `db:"payload"`
	SentAt  time.Time `db:"sent_at"`
}

func ToModel(notification *notificationentity.Notification) *Notification {
	return &Notification{
		ID:      notification.ID().String(),
		EventID: notification.EventID().String(),
		UserID:  notification.UserID().String(),
		Kind:    notification.Kind().String(),
		Payload: notification.Payload().Bytes(),
		SentAt:  notification.SentAt().Time(),
	}
}

func ToEntity(model *Notification) (*notificationentity.Notification, error) {
	notificationID, err := notificationentity.NotificationIDFrom(model.ID)
	if err != nil {
		return nil, err
	}

	eventID, err := notificationentity.EventIDFrom(model.EventID)
	if err != nil {
		return nil, err
	}

	userID, err := commonuserentity.UserIDFrom(model.UserID)
	if err != nil {
		return nil, err
	}

	kind, err := notificationentity.KindFrom(model.Kind)
	if err != nil {
		return nil, err
	}

	payload, err := notificationentity.PayloadFrom(model.Payload)
	if err != nil {
		return nil, err
	}

	return notificationentity.NewBuilder().
		ID(notificationID).
		EventID(eventID).
		UserID(userID).
		Kind(kind).
		Payload(payload).
		SentAt(timeprimitive.TimestampFrom(model.SentAt)).
		Build()
}
