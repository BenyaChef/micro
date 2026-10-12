package eventusecasehandler

import (
	"context"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	repositoryinterface "github.com/BenyaChef/micro/notification/boundary/repository"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
)

var _ Handler = (*RegisteredHandler)(nil)

type RegisteredHandler struct {
	notificationRepository repositoryinterface.NotificationRepository
	logPublisher           loggerinterface.LogPublisher
	errorProcessor         *errorProcessor
}

func (h *RegisteredHandler) Handle(ctx context.Context, event *commonautheventmodel.UserEvent) error {
	if event == nil {
		return h.errorProcessor.LogAndReturn(ctx, ErrEventIsRequired)
	}

	notification, err := h.buildNotification(event)
	if err != nil {
		return h.errorProcessor.LogAndReturn(ctx, err)
	}

	inserted, err := h.notificationRepository.Insert(ctx, notification)
	if err != nil {
		return h.errorProcessor.LogAndReturn(ctx, ErrStoreFailed(event.EventID, err))
	}

	if !inserted {
		h.logPublisher.LogInfo(ctx, "event already handled, skipping",
			"event_id", event.EventID,
			"user_id", event.Data.UserID,
		)

		return nil
	}

	h.logPublisher.LogInfo(ctx, "welcome notification stored",
		"event_id", event.EventID,
		"user_id", event.Data.UserID,
		"email", event.Data.Email,
	)

	return nil
}

func (h *RegisteredHandler) buildNotification(
	event *commonautheventmodel.UserEvent,
) (*notificationentity.Notification, error) {
	eventID, err := notificationentity.EventIDFrom(event.EventID)
	if err != nil {
		return nil, err
	}

	userID, err := commonuserentity.UserIDFrom(event.Data.UserID)
	if err != nil {
		return nil, err
	}

	rawPayload, err := event.ToBytes()
	if err != nil {
		return nil, err
	}

	payload, err := notificationentity.PayloadFrom(rawPayload)
	if err != nil {
		return nil, err
	}

	return notificationentity.NewBuilder().
		ID(notificationentity.NewNotificationID()).
		EventID(eventID).
		UserID(userID).
		Kind(notificationentity.KindWelcome).
		Payload(payload).
		Build()
}
