package postgres

import (
	"fmt"

	"micro/infrastructure/errors"
)

var ErrDSNIsRequired = errors.NewError("SYS", "Postgres: DSN is required")

const (
	ErrCodeInvalidDSN    = errors.ErrorCode("48db1e31-001")
	ErrCodeConnectFailed = errors.ErrorCode("48db1e31-040")
	ErrCodePingFailed    = errors.ErrorCode("48db1e31-041")
)

func ErrInvalidDSN(cause error) error {
	msg := "Postgres DSN is malformed"

	return errors.NewErrorWithLevelAndCause(ErrCodeInvalidDSN, msg, errors.LevelCritical, cause)
}

func ErrConnectFailed(host string, cause error) error {
	msg := fmt.Sprintf("Postgres connection failed: host=%s", host)

	return errors.NewErrorWithLevelAndCause(ErrCodeConnectFailed, msg, errors.LevelCritical, cause)
}

func ErrPingFailed(host string, cause error) error {
	msg := fmt.Sprintf("Postgres ping failed: host=%s", host)

	return errors.NewErrorWithLevelAndCause(ErrCodePingFailed, msg, errors.LevelCritical, cause)
}
