package idprimitive

import "micro/infrastructure/errors"

var (
	ErrEntityIDIsEmpty       = errors.NewErrorWithLevel("fed17d16-001", "Entity id is empty", errors.LevelInfo)
	ErrEntityIDInvalidFormat = errors.NewErrorWithLevel("fed17d16-002", "Entity id has invalid format", errors.LevelInfo)
)
