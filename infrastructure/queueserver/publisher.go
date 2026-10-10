package queueserver

import (
	"context"
	"strings"

	"github.com/IBM/sarama"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	queueserverinterface "github.com/BenyaChef/micro/infrastructure/queueserver/interface"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

var _ queueserverinterface.QueuePublisher = (*Publisher)(nil)

type Publisher struct {
	producer     sarama.SyncProducer
	logPublisher loggerinterface.LogPublisher
	brokers      []string
}

func (p *Publisher) SendMessage(ctx context.Context, message *queueservermodel.ProducerMessage) error {
	if message == nil {
		return ErrMessageIsRequired
	}

	partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: message.Topic,
		Key:   sarama.StringEncoder(message.PartitionKey),
		Value: sarama.ByteEncoder(message.Value),
	})
	if err != nil {
		return ErrSendFailed(message.Topic, err)
	}

	p.logPublisher.LogInfo(ctx, "message published",
		"topic", message.Topic,
		"key", message.PartitionKey,
		"partition", partition,
		"offset", offset,
	)

	return nil
}

func (p *Publisher) Brokers() string {
	return strings.Join(p.brokers, ",")
}

func (p *Publisher) Close(ctx context.Context) {
	if p.producer == nil {
		return
	}

	if err := p.producer.Close(); err != nil {
		p.logPublisher.LogError(ctx, err)
	}
}
