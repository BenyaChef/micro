package queueserver

import (
	"context"
	"errors"
	"strings"

	"github.com/IBM/sarama"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	queueserverinterface "github.com/BenyaChef/micro/infrastructure/queueserver/interface"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

var _ sarama.ConsumerGroupHandler = (*Consumer)(nil)

type Consumer struct {
	group        sarama.ConsumerGroup
	handler      queueserverinterface.QueueHandler
	logPublisher loggerinterface.LogPublisher
	topics       []string
	groupID      string
	semaphore    chan struct{}
}

func (c *Consumer) Run(ctx context.Context) error {
	c.logPublisher.LogInfo(ctx, "queue consumer started",
		"topics", strings.Join(c.topics, ","),
		"group", c.groupID,
	)

	for ctx.Err() == nil {
		err := c.group.Consume(ctx, c.topics, c)
		if err == nil {
			continue
		}

		if isShutdown(ctx, err) {
			break
		}

		return ErrConsumeFailed(strings.Join(c.topics, ","), err)
	}

	c.logPublisher.LogInfo(ctx, "queue consumer stopped", "group", c.groupID)

	return nil
}

func isShutdown(ctx context.Context, err error) bool {
	return ctx.Err() != nil ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, sarama.ErrClosedConsumerGroup)
}

func (c *Consumer) Close(ctx context.Context) {
	if c.group == nil {
		return
	}

	if err := c.group.Close(); err != nil {
		c.logPublisher.LogError(ctx, err)
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()

	for saramaMessage := range claim.Messages() {
		if err := c.handle(ctx, saramaMessage); err != nil {
			return err
		}

		session.MarkMessage(saramaMessage, "")
	}

	return nil
}

func (c *Consumer) handle(ctx context.Context, saramaMessage *sarama.ConsumerMessage) error {
	select {
	case c.semaphore <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}

	defer func() { <-c.semaphore }()

	message := &queueservermodel.ConsumerMessage{
		Topic:     saramaMessage.Topic,
		Key:       string(saramaMessage.Key),
		Value:     saramaMessage.Value,
		Partition: saramaMessage.Partition,
		Offset:    saramaMessage.Offset,
	}

	return c.handler.Handle(ctx, message)
}
