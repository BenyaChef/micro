package queueserverinterface

import (
	"context"

	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

type QueueHandler interface {
	Handle(ctx context.Context, message *queueservermodel.ConsumerMessage) error
}
