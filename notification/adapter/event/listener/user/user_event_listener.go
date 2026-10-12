package usereventlistener

import (
	"context"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	queueserverinterface "github.com/BenyaChef/micro/infrastructure/queueserver/interface"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
	usecaseinterface "github.com/BenyaChef/micro/notification/boundary/usecase"
)

var _ queueserverinterface.QueueHandler = (*UserEventListener)(nil)

type UserEventListener struct {
	userEventUseCase usecaseinterface.UserEventUseCase
	logPublisher     loggerinterface.LogPublisher
}

func (l *UserEventListener) Handle(ctx context.Context, message *queueservermodel.ConsumerMessage) error {
	if message == nil {
		l.logPublisher.LogError(ctx, ErrMessageIsRequired)

		return nil
	}

	event, err := commonautheventmodel.UserEventFromBytes(message.Value)
	if err != nil {
		l.logPublisher.LogWarn(ctx, "message cannot be parsed, skipping",
			"topic", message.Topic,
			"partition", message.Partition,
			"offset", message.Offset,
			"cause", err.Error(),
		)

		return nil
	}

	return l.userEventUseCase.Handle(ctx, event)
}
