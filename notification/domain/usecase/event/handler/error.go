package eventusecasehandler

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const ErrCodeStoreFailed = errors.ErrorCode("c1195890-010")

var (
	ErrNotificationRepositoryIsRequired = errors.NewError(
		"SYS", "RegisteredHandler: NotificationRepository is required",
	)
	ErrLogPublisherIsRequired = errors.NewError("SYS", "RegisteredHandler: LogPublisher is required")
	ErrHandlersAreRequired    = errors.NewError("SYS", "HandlerProvider: at least one handler is required")
)

var ErrEventIsRequired = errors.NewErrorWithLevel(
	"c1195890-001", "RegisteredHandler: event is required", errors.LevelInfo,
)

func ErrStoreFailed(eventID string, cause error) error {
	msg := fmt.Sprintf("Unable to store notification: eventID=%s", eventID)

	return errors.NewErrorWithCause(ErrCodeStoreFailed, msg, cause)
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
