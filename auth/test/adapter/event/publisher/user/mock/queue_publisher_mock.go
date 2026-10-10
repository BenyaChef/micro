package usereventpublishermock

import (
	"context"

	queueserverinterface "github.com/BenyaChef/micro/infrastructure/queueserver/interface"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

var _ queueserverinterface.QueuePublisher = (*QueuePublisherMock)(nil)

type QueuePublisherMock struct {
	SendMessageFunc func(ctx context.Context, message *queueservermodel.ProducerMessage) error

	Sent []*queueservermodel.ProducerMessage
}

func NewQueuePublisherMock() *QueuePublisherMock {
	return &QueuePublisherMock{
		SendMessageFunc: func(context.Context, *queueservermodel.ProducerMessage) error { return nil },
	}
}

func (m *QueuePublisherMock) SendMessage(ctx context.Context, message *queueservermodel.ProducerMessage) error {
	m.Sent = append(m.Sent, message)

	return m.SendMessageFunc(ctx, message)
}
