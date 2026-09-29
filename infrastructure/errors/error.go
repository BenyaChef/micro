package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	code  ErrorCode
	msg   string
	level ErrorLevel
	cause error
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func (e *Error) Message() string {
	return e.msg
}

func (e *Error) Level() ErrorLevel {
	return e.level
}

func (e *Error) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("%s: %s. Level - '%s'", e.code, e.msg, e.level)
	}

	return fmt.Sprintf("%s: %s. Level - '%s'. Cause - %v", e.code, e.msg, e.level, e.cause)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	var targetErr *Error
	if errors.As(target, &targetErr) {
		return e.code == targetErr.code && e.msg == targetErr.msg
	}

	return false
}
