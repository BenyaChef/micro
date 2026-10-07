package postentity

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeTitleTooLong = errors.ErrorCode("59831142-005")
	ErrCodeBodyTooLong  = errors.ErrorCode("59831142-006")
)

var (
	ErrPostIDInvalidFormat = errors.NewErrorWithLevel(
		"59831142-001", "Post id has invalid format", errors.LevelInfo,
	)
	ErrTitleIsEmpty = errors.NewErrorWithLevel(
		"59831142-002", "Post title is empty", errors.LevelInfo,
	)
	ErrBodyIsEmpty = errors.NewErrorWithLevel(
		"59831142-003", "Post body is empty", errors.LevelInfo,
	)
)

var (
	ErrPostIDIsRequired   = errors.NewError("59831142-010", "Post: PostID is required")
	ErrAuthorIDIsRequired = errors.NewError("59831142-011", "Post: AuthorID is required")
	ErrTitleIsRequired    = errors.NewError("59831142-012", "Post: Title is required")
	ErrBodyIsRequired     = errors.NewError("59831142-013", "Post: Body is required")
)

func ErrTitleTooLong(length, maxLength int) error {
	msg := fmt.Sprintf("Post title is too long: length=%d, max=%d", length, maxLength)

	return errors.NewErrorWithLevel(ErrCodeTitleTooLong, msg, errors.LevelInfo)
}

func ErrBodyTooLong(length, maxLength int) error {
	msg := fmt.Sprintf("Post body is too long: length=%d, max=%d", length, maxLength)

	return errors.NewErrorWithLevel(ErrCodeBodyTooLong, msg, errors.LevelInfo)
}
