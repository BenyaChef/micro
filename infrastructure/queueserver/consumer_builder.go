package queueserver

import (
	"errors"
	"strings"

	"github.com/IBM/sarama"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	queueserverinterface "github.com/BenyaChef/micro/infrastructure/queueserver/interface"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

type ConsumerBuilder struct {
	brokers        []string
	topics         []string
	groupID        string
	handler        queueserverinterface.QueueHandler
	logPublisher   loggerinterface.LogPublisher
	maxConcurrency int
}

func NewConsumerBuilder() *ConsumerBuilder {
	return &ConsumerBuilder{maxConcurrency: queueservermodel.DefaultMaxConcurrentHandlers}
}

func (b *ConsumerBuilder) Brokers(brokers string) *ConsumerBuilder {
	b.brokers = splitAndTrim(brokers)

	return b
}

func (b *ConsumerBuilder) Topics(topics ...string) *ConsumerBuilder {
	for _, topic := range topics {
		if trimmed := strings.TrimSpace(topic); trimmed != "" {
			b.topics = append(b.topics, trimmed)
		}
	}

	return b
}

func (b *ConsumerBuilder) GroupID(groupID string) *ConsumerBuilder {
	b.groupID = groupID

	return b
}

func (b *ConsumerBuilder) Handler(handler queueserverinterface.QueueHandler) *ConsumerBuilder {
	b.handler = handler

	return b
}

func (b *ConsumerBuilder) LogPublisher(logPublisher loggerinterface.LogPublisher) *ConsumerBuilder {
	b.logPublisher = logPublisher

	return b
}

func (b *ConsumerBuilder) MaxConcurrency(maxConcurrency int) *ConsumerBuilder {
	b.maxConcurrency = maxConcurrency

	return b
}

func (b *ConsumerBuilder) Build() (*Consumer, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(b.brokers, b.groupID, config)
	if err != nil {
		return nil, ErrConsumerGroupInit(b.groupID, err)
	}

	return &Consumer{
		group:        group,
		handler:      b.handler,
		logPublisher: b.logPublisher,
		topics:       b.topics,
		groupID:      b.groupID,
		semaphore:    make(chan struct{}, b.maxConcurrency),
	}, nil
}

func (b *ConsumerBuilder) checkRequiredFields() error {
	var errs []error

	if len(b.brokers) == 0 {
		errs = append(errs, ErrConsumerBrokersAreRequired)
	}

	if len(b.topics) == 0 {
		errs = append(errs, ErrConsumerTopicsAreRequired)
	}

	if b.groupID == "" {
		errs = append(errs, ErrConsumerGroupIDIsRequired)
	}

	if b.handler == nil {
		errs = append(errs, ErrConsumerHandlerIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrConsumerLogPublisherIsRequired)
	}

	if b.maxConcurrency <= 0 {
		b.maxConcurrency = queueservermodel.DefaultMaxConcurrentHandlers
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func splitAndTrim(raw string) []string {
	var values []string

	for _, value := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			values = append(values, trimmed)
		}
	}

	return values
}
