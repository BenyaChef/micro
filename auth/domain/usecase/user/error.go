package userusecase

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	ErrCodePasswordTooShort  = errors.ErrorCode("d3c2ee1e-003")
	ErrCodeOperationFailed   = errors.ErrorCode("d3c2ee1e-010")
	ErrCodeUserNotFound      = errors.ErrorCode("d3c2ee1e-020")
	ErrCodeEmailAlreadyTaken = errors.ErrorCode("d3c2ee1e-030")
)

var (
	ErrUserRepositoryIsRequired = errors.NewError("SYS", "UserUseCase: UserRepository is required")
	ErrTokenUseCaseIsRequired   = errors.NewError("SYS", "UserUseCase: TokenUseCase is required")
	ErrPasswordHasherIsRequired = errors.NewError("SYS", "UserUseCase: PasswordHasher is required")
	ErrLogPublisherIsRequired   = errors.NewError("SYS", "UserUseCase: LogPublisher is required")
)

var (
	ErrDataIsRequired     = errors.NewErrorWithLevel("d3c2ee1e-001", "User: request data is required", errors.LevelInfo)
	ErrPasswordIsRequired = errors.NewErrorWithLevel("d3c2ee1e-002", "User: password is required", errors.LevelInfo)
	ErrInvalidCredentials = errors.NewErrorWithLevel("d3c2ee1e-031", "Invalid email or password", errors.LevelInfo)
)

func ErrPasswordTooShort(length, minLength int) error {
	msg := fmt.Sprintf("Password is too short: length=%d, min=%d", length, minLength)

	return errors.NewErrorWithLevel(ErrCodePasswordTooShort, msg, errors.LevelInfo)
}

func ErrUserNotFound(userID string) error {
	msg := fmt.Sprintf("User not found: userID=%s", userID)

	return errors.NewErrorWithLevel(ErrCodeUserNotFound, msg, errors.LevelInfo)
}

func ErrEmailAlreadyTaken(email string) error {
	msg := fmt.Sprintf("Email is already taken: email=%s", email)

	return errors.NewErrorWithLevel(ErrCodeEmailAlreadyTaken, msg, errors.LevelInfo)
}

func ErrOperationFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeOperationFailed, "User operation failed", cause)
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

func (p *errorProcessor) LogAndReturnOperationError(ctx context.Context, cause error) error {
	err := ErrOperationFailed(cause)
	p.logPublisher.LogError(ctx, err)

	return err
}
