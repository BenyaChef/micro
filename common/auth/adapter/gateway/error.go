package commonauthgateway

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const ErrCodeCallFailed = errors.ErrorCode("52dcc9fb-040")

var (
	ErrClientIsRequired       = errors.NewError("SYS", "AuthGateway: Client is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "AuthGateway: LogPublisher is required")
)

var (
	ErrTokenIsRequired = errors.NewErrorWithLevel(
		"52dcc9fb-001", "AuthGateway: raw token is required", errors.LevelInfo,
	)
	ErrEmptyResult = errors.NewError(
		"52dcc9fb-002", "AuthGateway: response carries neither token nor error",
	)
)

func ErrCallFailed(cause error) error {
	msg := fmt.Sprintf("Unable to call auth service: method=%s", "ValidateToken")

	return errors.NewErrorWithCause(ErrCodeCallFailed, msg, cause)
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

func (p *errorProcessor) LogAndReturnCallError(ctx context.Context, cause error) error {
	err := ErrCallFailed(cause)
	p.logPublisher.LogError(ctx, err)

	return err
}
