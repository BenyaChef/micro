package usereventpublisher

import (
	"errors"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	queueserverinterface "github.com/BenyaChef/micro/infrastructure/queueserver/interface"
)

type Builder struct {
	queuePublisher queueserverinterface.QueuePublisher
	logPublisher   loggerinterface.LogPublisher
	topic          string
}

func NewBuilder() *Builder {
	return &Builder{topic: commonautheventmodel.TopicUser}
}

func (b *Builder) QueuePublisher(queuePublisher queueserverinterface.QueuePublisher) *Builder {
	b.queuePublisher = queuePublisher

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Topic(topic string) *Builder {
	b.topic = topic

	return b
}

func (b *Builder) Build() (*UserEventPublisher, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	return &UserEventPublisher{
		queuePublisher: b.queuePublisher,
		logPublisher:   b.logPublisher,
		topic:          b.topic,
		errorProcessor: newErrorProcessor(b.logPublisher),
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.queuePublisher == nil {
		errs = append(errs, ErrQueuePublisherIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if b.topic == "" {
		errs = append(errs, ErrTopicIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
