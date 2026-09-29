package errors

func NewError(code ErrorCode, msg string) *Error {
	return &Error{code: code, msg: msg, level: LevelError}
}

func NewErrorWithLevel(code ErrorCode, msg string, level ErrorLevel) *Error {
	return &Error{code: code, msg: msg, level: level}
}

func NewErrorWithCause(code ErrorCode, msg string, cause error) *Error {
	return &Error{code: code, msg: msg, level: LevelError, cause: cause}
}

func NewErrorWithLevelAndCause(code ErrorCode, msg string, level ErrorLevel, cause error) *Error {
	return &Error{code: code, msg: msg, level: level, cause: cause}
}

func NewErrorFrom(err error) *Error {
	return CastOrWrap(err, UnknownErrorCode)
}
