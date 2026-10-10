package queueserver

import (
	"errors"
	"strings"
	"time"

	"github.com/IBM/sarama"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

type Builder struct {
	brokers      []string
	logPublisher loggerinterface.LogPublisher
	retryMax     int
	timeout      time.Duration
}

func NewBuilder() *Builder {
	return &Builder{
		retryMax: queueservermodel.DefaultRetryMax,
		timeout:  queueservermodel.DefaultTimeout * time.Second,
	}
}

func (b *Builder) Brokers(brokers string) *Builder {
	for _, broker := range strings.Split(brokers, ",") {
		if trimmed := strings.TrimSpace(broker); trimmed != "" {
			b.brokers = append(b.brokers, trimmed)
		}
	}

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) RetryMax(retryMax int) *Builder {
	b.retryMax = retryMax

	return b
}

func (b *Builder) Timeout(timeout time.Duration) *Builder {
	b.timeout = timeout

	return b
}

func (b *Builder) Build() (*Publisher, error) {
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

	return &Publisher{producer: producer, logPublisher: b.logPublisher, brokers: b.brokers}, nil
}

func (b *Builder) checkRequiredFields() error {
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
