package initservice

import (
	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	eventusecase "github.com/BenyaChef/micro/notification/domain/usecase/event"
	eventusecasehandler "github.com/BenyaChef/micro/notification/domain/usecase/event/handler"
)

func (c *DependencyContainer) initUseCases() error {
	registeredHandler, err := eventusecasehandler.NewBuilder().
		NotificationRepository(c.notificationRepository).
		LogPublisher(c.infra.LogPublisher()).
		BuildRegisteredHandler()
	if err != nil {
		return err
	}

	provider, err := eventusecasehandler.NewProvider(map[commonautheventmodel.UserEventAction]eventusecasehandler.Handler{
		commonautheventmodel.UserEventActionRegistered: registeredHandler,
	})
	if err != nil {
		return err
	}

	userEventUseCase, err := eventusecase.NewBuilder().
		HandlerProvider(provider).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.userEventUseCase = userEventUseCase

	return nil
}
