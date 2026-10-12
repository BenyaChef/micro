package notificationrepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	postgresinterface "github.com/BenyaChef/micro/infrastructure/postgres/interface"
	notificationrepomodel "github.com/BenyaChef/micro/notification/adapter/repository/notification/model"
	repositoryinterface "github.com/BenyaChef/micro/notification/boundary/repository"
	notificationentity "github.com/BenyaChef/micro/notification/domain/entity/notification"
)

const (
	queryInsert = `
		INSERT INTO notifications (id, event_id, user_id, kind, payload, sent_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (event_id) DO NOTHING`

	queryGetByEventID = `
		SELECT id, event_id, user_id, kind, payload, sent_at
		FROM notifications
		WHERE event_id = $1`
)

var _ repositoryinterface.NotificationRepository = (*NotificationRepository)(nil)

type NotificationRepository struct {
	executor       postgresinterface.Executor
	logPublisher   loggerinterface.LogPublisher
	errorProcessor *errorProcessor
}

func (r *NotificationRepository) Insert(
	ctx context.Context,
	notification *notificationentity.Notification,
) (bool, error) {
	if notification == nil {
		return false, r.errorProcessor.LogAndReturn(ctx, ErrNotificationIsRequired)
	}

	model := notificationrepomodel.ToModel(notification)

	tag, err := r.executor.Exec(ctx, queryInsert,
		model.ID, model.EventID, model.UserID, model.Kind, model.Payload, model.SentAt,
	)
	if err != nil {
		return false, r.errorProcessor.LogAndReturn(ctx, ErrInsertFailed(model.EventID, err))
	}

	return tag.RowsAffected() > 0, nil
}

func (r *NotificationRepository) GetByEventID(
	ctx context.Context,
	eventID notificationentity.EventID,
) (*notificationentity.Notification, error) {
	rows, err := r.executor.Query(ctx, queryGetByEventID, eventID.String())
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	model, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[notificationrepomodel.Notification])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	notification, err := notificationrepomodel.ToEntity(model)
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, err)
	}

	return notification, nil
}
