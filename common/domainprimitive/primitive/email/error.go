package emailprimitive

import "github.com/BenyaChef/micro/infrastructure/errors"

var (
	ErrEmailIsEmpty       = errors.NewErrorWithLevel("545af577-001", "Email is empty", errors.LevelInfo)
	ErrEmailInvalidFormat = errors.NewErrorWithLevel("545af577-002", "Email has invalid format", errors.LevelInfo)
)
