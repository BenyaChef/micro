package userrepository

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

var (
	ErrExecutorIsRequired     = errors.NewError("SYS", "UserRepository: Executor is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "UserRepository: LogPublisher is required")
)

const (
	ErrCodeInsertFailed  = errors.ErrorCode("fca1b026-010")
	ErrCodeSelectFailed  = errors.ErrorCode("fca1b026-011")
	ErrCodeEmailConflict = errors.ErrorCode("fca1b026-030")
)

var ErrUserIsRequired = errors.NewErrorWithLevel(
	"fca1b026-001", "UserRepository: user is required", errors.LevelInfo,
)

func ErrInsertFailed(userID string, cause error) error {
	msg := fmt.Sprintf("Unable to insert user: userID=%s", userID)

	return errors.NewErrorWithCause(ErrCodeInsertFailed, msg, cause)
}

func ErrSelectFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeSelectFailed, "Unable to select user", cause)
}

func ErrEmailConflict(email string, cause error) error {
	msg := fmt.Sprintf("User with this email already exists: email=%s", email)

	return errors.NewErrorWithLevelAndCause(ErrCodeEmailConflict, msg, errors.LevelInfo, cause)
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
