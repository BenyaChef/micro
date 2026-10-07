package authgateway

import (
	"context"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	ErrCodeTokenRejected   = errors.ErrorCode("2564e485-030")
	ErrCodeAuthUnavailable = errors.ErrorCode("2564e485-040")
)

var (
	ErrAuthGatewayIsRequired  = errors.NewError("SYS", "BlogAuthGateway: AuthGateway is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "BlogAuthGateway: LogPublisher is required")
)

func ErrTokenRejected(cause error) error {
	return errors.NewErrorWithLevelAndCause(
		ErrCodeTokenRejected, "Access token was rejected by auth service", errors.LevelInfo, cause,
	)
}

func ErrAuthUnavailable(cause error) error {
	return errors.NewErrorWithLevelAndCause(
		ErrCodeAuthUnavailable, "Auth service is unavailable", errors.LevelError, cause,
	)
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
