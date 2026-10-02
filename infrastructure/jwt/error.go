package jwt

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeSecretTooShort = errors.ErrorCode("10b754c6-001")
	ErrCodeTokenInvalid   = errors.ErrorCode("10b754c6-002")
	ErrCodeTokenExpired   = errors.ErrorCode("10b754c6-003")
	ErrCodeSubjectIsEmpty = errors.ErrorCode("10b754c6-004")
	ErrCodeSignFailed     = errors.ErrorCode("10b754c6-010")
)

var (
	ErrSecretIsRequired = errors.NewError("SYS", "JWT: Secret is required")
	ErrIssuerIsRequired = errors.NewError("SYS", "JWT: Issuer is required")
)

var (
	ErrTokenInvalid = errors.NewErrorWithLevel(
		ErrCodeTokenInvalid, "Token is invalid", errors.LevelInfo,
	)
	ErrTokenExpired = errors.NewErrorWithLevel(
		ErrCodeTokenExpired, "Token is expired", errors.LevelInfo,
	)
	ErrSubjectIsEmpty = errors.NewErrorWithLevel(
		ErrCodeSubjectIsEmpty, "Token subject is empty", errors.LevelInfo,
	)
)

func ErrSecretTooShort(length, minLength int) error {
	msg := fmt.Sprintf("JWT secret is too short: length=%d, min=%d", length, minLength)

	return errors.NewErrorWithLevel(ErrCodeSecretTooShort, msg, errors.LevelCritical)
}

func ErrSignFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeSignFailed, "Unable to sign token", cause)
}
