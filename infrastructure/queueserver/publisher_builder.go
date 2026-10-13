package queueserver

import (
	"errors"
	"strings"
	"time"

	"github.com/IBM/sarama"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	infraotel "github.com/BenyaChef/micro/infrastructure/otel"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

type PublisherBuilder struct {
	brokers      []string
	logPublisher loggerinterface.LogPublisher
	kafkaMonitor *infraotel.KafkaMonitor
	retryMax     int
	timeout      time.Duration
}

func NewPublisherBuilder() *PublisherBuilder {
	return &PublisherBuilder{
		retryMax: queueservermodel.DefaultRetryMax,
		timeout:  queueservermodel.DefaultTimeout * time.Second,
	}
}

func (b *PublisherBuilder) Brokers(brokers string) *PublisherBuilder {
	b.brokers = splitAndTrim(brokers)

	return b
}

func (b *PublisherBuilder) LogPublisher(logPublisher loggerinterface.LogPublisher) *PublisherBuilder {
	b.logPublisher = logPublisher

	return b
}

func (b *PublisherBuilder) KafkaMonitor(kafkaMonitor *infraotel.KafkaMonitor) *PublisherBuilder {
	b.kafkaMonitor = kafkaMonitor

	return b
}

func (b *PublisherBuilder) RetryMax(retryMax int) *PublisherBuilder {
	b.retryMax = retryMax

	return b
}

func (b *PublisherBuilder) Timeout(timeout time.Duration) *PublisherBuilder {
	b.timeout = timeout

	return b
}

func (b *PublisherBuilder) Build() (*Publisher, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = b.retryMax
	config.Producer.Return.Successes = true
	config.Producer.Timeout = b.timeout

	producer, err := sarama.NewSyncProducer(b.brokers, config)
	if err != nil {
		return nil, ErrConnectFailed(strings.Join(b.brokers, ","), err)
	}

	return &Publisher{
		producer:     producer,
		logPublisher: b.logPublisher,
		kafkaMonitor: b.kafkaMonitor,
		brokers:      b.brokers,
	}, nil
}

func (b *PublisherBuilder) checkRequiredFields() error {
	var errs []error

	if len(b.brokers) == 0 {
		errs = append(errs, ErrBrokersAreRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
