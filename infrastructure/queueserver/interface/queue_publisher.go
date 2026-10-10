package queueserverinterface

import (
	"context"

	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

type QueuePublisher interface {
	SendMessage(ctx context.Context, message *queueservermodel.ProducerMessage) error
}
