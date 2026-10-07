package restservermodel

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeUnmarshalRequest      = errors.ErrorCode("64c4a693-001")
	ErrCodeMissingRequiredField  = errors.ErrorCode("64c4a693-002")
	ErrCodeAuthorizationRequired = errors.ErrorCode("64c4a693-003")
	ErrCodeWriteResponse         = errors.ErrorCode("64c4a693-010")
)

var (
	ErrWriteResponse = errors.NewError(ErrCodeWriteResponse, "Failed to write response")

	ErrAuthorizationRequired = errors.NewErrorWithLevel(
		ErrCodeAuthorizationRequired, "Authorization header with Bearer token is required", errors.LevelInfo,
	)
)

func ErrUnmarshalRequest(cause string) error {
	msg := fmt.Sprintf("Failed to parse request body: cause=%s", cause)

	return errors.NewErrorWithLevel(ErrCodeUnmarshalRequest, msg, errors.LevelInfo)
}

func ErrMissingRequiredField(field string) error {
	msg := fmt.Sprintf("Missing required field in request: field=%s", field)

	return errors.NewErrorWithLevel(ErrCodeMissingRequiredField, msg, errors.LevelInfo)
}
