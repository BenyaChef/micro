package notificationrepository

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	ErrCodeInsertFailed = errors.ErrorCode("6ce2099a-010")
	ErrCodeSelectFailed = errors.ErrorCode("6ce2099a-011")
)

var (
	ErrExecutorIsRequired     = errors.NewError("SYS", "NotificationRepository: Executor is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "NotificationRepository: LogPublisher is required")
)

var ErrNotificationIsRequired = errors.NewErrorWithLevel(
	"6ce2099a-001", "NotificationRepository: notification is required", errors.LevelInfo,
)

func ErrInsertFailed(eventID string, cause error) error {
	msg := fmt.Sprintf("Unable to insert notification: eventID=%s", eventID)

	return errors.NewErrorWithCause(ErrCodeInsertFailed, msg, cause)
}

func ErrSelectFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeSelectFailed, "Unable to select notification", cause)
}

type errorProcessor struct {
	logPublisher loggerinterface.LogPublisher
}

func newErrorProcessor(logPublisher loggerinterface.LogPublisher) *errorProcessor {
	return &errorProcessor{logPublisher: logPublisher}
}

func (p *errorProcessor) LogAndReturn(ctx context.Context, err error) error {
	p.logPublisher.LogError(ctx, err)

	return err
}
