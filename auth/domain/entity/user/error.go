package userentity

import "github.com/BenyaChef/micro/infrastructure/errors"

var (
	ErrPasswordHashIsEmpty = errors.NewErrorWithLevel(
		"590ac348-001", "Password hash is empty", errors.LevelInfo,
	)
	ErrUserIDIsRequired = errors.NewError(
		"590ac348-002", "User: UserID is required",
	)
	ErrEmailIsRequired = errors.NewError(
		"590ac348-003", "User: Email is required",
	)
	ErrPasswordHashIsRequired = errors.NewError(
		"590ac348-004", "User: PasswordHash is required",
	)
)
