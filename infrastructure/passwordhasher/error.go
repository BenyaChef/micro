package passwordhasher

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodePasswordIsEmpty    = errors.ErrorCode("ce0a47e0-001")
	ErrCodeCostOutOfRange     = errors.ErrorCode("ce0a47e0-002")
	ErrCodePasswordMismatch   = errors.ErrorCode("ce0a47e0-003")
	ErrCodePasswordHashFailed = errors.ErrorCode("ce0a47e0-010")
)

var (
	ErrPasswordIsEmpty = errors.NewErrorWithLevel(
		ErrCodePasswordIsEmpty, "Password is empty", errors.LevelInfo,
	)
	ErrPasswordMismatch = errors.NewErrorWithLevel(
		ErrCodePasswordMismatch, "Password does not match", errors.LevelInfo,
	)
)

func ErrCostOutOfRange(cost, minCost, maxCost int) error {
	msg := fmt.Sprintf("Password hash cost is out of range: cost=%d, min=%d, max=%d", cost, minCost, maxCost)

	return errors.NewErrorWithLevel(ErrCodeCostOutOfRange, msg, errors.LevelCritical)
}

func ErrPasswordHashFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodePasswordHashFailed, "Unable to hash password", cause)
}
