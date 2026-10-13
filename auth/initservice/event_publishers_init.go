package initservice

import (
	usereventpublisher "github.com/BenyaChef/micro/auth/adapter/event/publisher/user"
	"github.com/BenyaChef/micro/infrastructure/queueserver"
)

func (c *DependencyContainer) initQueuePublisher() error {
	publisher, err := queueserver.NewPublisherBuilder().
		Brokers(c.environments.kafkaBrokers).
		LogPublisher(c.infra.LogPublisher()).
		KafkaMonitor(c.infra.KafkaMonitor()).
		Build()
	if err != nil {
		return err
	}

	c.queuePublisher = publisher

	return nil
}

func (c *DependencyContainer) initEventPublishers() error {
	publisher, err := usereventpublisher.NewBuilder().
		QueuePublisher(c.queuePublisher).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.userEventPublisher = publisher

	return nil
}
