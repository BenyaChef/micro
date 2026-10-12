package eventusecase

import (
	"context"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	usecaseinterface "github.com/BenyaChef/micro/notification/boundary/usecase"
	eventusecasehandler "github.com/BenyaChef/micro/notification/domain/usecase/event/handler"
)

var _ usecaseinterface.UserEventUseCase = (*UserEventUseCase)(nil)

type UserEventUseCase struct {
	handlerProvider *eventusecasehandler.Provider
	logPublisher    loggerinterface.LogPublisher
	errorProcessor  *errorProcessor
}

func (uc *UserEventUseCase) Handle(ctx context.Context, event *commonautheventmodel.UserEvent) error {
	if event == nil {
		return uc.errorProcessor.LogAndReturn(ctx, ErrEventIsRequired)
	}

	if err := event.Validate(); err != nil {
		return uc.errorProcessor.LogAndReturn(ctx, err)
	}

	handler, found := uc.handlerProvider.ByAction(event.Action)
	if !found {
		uc.logPublisher.LogWarn(ctx, "no handler for action, skipping",
			"action", event.Action.String(),
			"event_id", event.EventID,
		)

		return nil
	}

	return handler.Handle(ctx, event)
}
