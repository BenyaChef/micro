package usereventpublisher

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	ErrCodeMarshalFailed = errors.ErrorCode("f3b60a0b-010")
	ErrCodePublishFailed = errors.ErrorCode("f3b60a0b-040")
)

var (
	ErrQueuePublisherIsRequired = errors.NewError("SYS", "UserEventPublisher: QueuePublisher is required")
	ErrLogPublisherIsRequired   = errors.NewError("SYS", "UserEventPublisher: LogPublisher is required")
	ErrTopicIsRequired          = errors.NewError("SYS", "UserEventPublisher: Topic is required")
)

var ErrUserIDIsRequired = errors.NewErrorWithLevel(
	"f3b60a0b-001", "UserEventPublisher: userID is required", errors.LevelInfo,
)

func ErrMarshalFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeMarshalFailed, "Unable to serialize user event", cause)
}

func ErrPublishFailed(userID string, cause error) error {
	msg := fmt.Sprintf("Unable to publish user event: userID=%s", userID)

	return errors.NewErrorWithCause(ErrCodePublishFailed, msg, cause)
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
