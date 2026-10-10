package usereventpublisher

import (
	"context"
	"time"

	eventinterface "github.com/BenyaChef/micro/auth/boundary/event"
	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	idprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/id"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	queueserverinterface "github.com/BenyaChef/micro/infrastructure/queueserver/interface"
	queueservermodel "github.com/BenyaChef/micro/infrastructure/queueserver/model"
)

var _ eventinterface.UserEventPublisher = (*UserEventPublisher)(nil)

type UserEventPublisher struct {
	queuePublisher queueserverinterface.QueuePublisher
	logPublisher   loggerinterface.LogPublisher
	topic          string
	errorProcessor *errorProcessor
}

func (p *UserEventPublisher) Registered(
	ctx context.Context,
	userID commonuserentity.UserID,
	email emailprimitive.Email,
) error {
	if userID.IsZero() {
		return p.errorProcessor.LogAndReturn(ctx, ErrUserIDIsRequired)
	}

	event := commonautheventmodel.NewUserRegistered(
		idprimitive.NewEntityID().String(),
		userID.String(),
		email.String(),
		time.Now(),
	)

	value, err := event.ToBytes()
	if err != nil {
		return p.errorProcessor.LogAndReturn(ctx, ErrMarshalFailed(err))
	}

	message := &queueservermodel.ProducerMessage{
		Topic:        p.topic,
		PartitionKey: event.PartitionKey(),
		Value:        value,
	}

	if err := p.queuePublisher.SendMessage(ctx, message); err != nil {
		return p.errorProcessor.LogAndReturn(ctx, ErrPublishFailed(userID.String(), err))
	}

	return nil
}
