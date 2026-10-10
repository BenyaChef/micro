package commonautheventmodel

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeActionUnsupported = errors.ErrorCode("298ef509-004")
	ErrCodeUnmarshalFailed   = errors.ErrorCode("298ef509-010")
)

var (
	ErrEventIDIsEmpty = errors.NewErrorWithLevel(
		"298ef509-001", "UserEvent: event_id is empty", errors.LevelInfo,
	)
	ErrActionIsEmpty = errors.NewErrorWithLevel(
		"298ef509-002", "UserEvent: action is empty", errors.LevelInfo,
	)
	ErrUserIDIsEmpty = errors.NewErrorWithLevel(
		"298ef509-003", "UserEvent: data.user_id is empty", errors.LevelInfo,
	)
)

func ErrActionUnsupported(rawAction string) error {
	msg := fmt.Sprintf("UserEvent: unsupported action: action=%s", rawAction)

	return errors.NewErrorWithLevel(ErrCodeActionUnsupported, msg, errors.LevelWarn)
}

func ErrUnmarshalFailed(cause error) error {
	return errors.NewErrorWithLevelAndCause(
		ErrCodeUnmarshalFailed, "UserEvent: unable to parse message", errors.LevelWarn, cause,
	)
}
