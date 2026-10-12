package initservice

import (
	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	"github.com/BenyaChef/micro/infrastructure/queueserver"
	usereventlistener "github.com/BenyaChef/micro/notification/adapter/event/listener/user"
)

func (c *DependencyContainer) initEventListeners() error {
	listener, err := usereventlistener.NewBuilder().
		UserEventUseCase(c.userEventUseCase).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	consumer, err := queueserver.NewConsumerBuilder().
		Brokers(c.environments.kafkaBrokers).
		Topics(commonautheventmodel.TopicUser).
		GroupID(c.environments.kafkaGroupID).
		Handler(listener).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.queueConsumer = consumer

	return nil
}
