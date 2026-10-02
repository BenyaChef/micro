package tokenusecase

import (
	"context"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const ErrCodeIssueFailed = errors.ErrorCode("9e65307e-010")

var (
	ErrTokenServiceIsRequired = errors.NewError("SYS", "TokenUseCase: TokenService is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "TokenUseCase: LogPublisher is required")
)

var (
	ErrUserIDIsRequired = errors.NewErrorWithLevel("9e65307e-001", "Token: UserID is required", errors.LevelInfo)
	ErrTokenIsRequired  = errors.NewErrorWithLevel("9e65307e-002", "Token: raw token is required", errors.LevelInfo)
)

func ErrIssueFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeIssueFailed, "Unable to issue access token", cause)
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

func (p *errorProcessor) LogAndReturnIssueError(ctx context.Context, cause error) error {
	err := ErrIssueFailed(cause)
	p.logPublisher.LogError(ctx, err)

	return err
}
