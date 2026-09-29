package errors

const UnknownErrorCode = ErrorCode("unknown_error")

type ErrorCode string

func (c ErrorCode) String() string {
	return string(c)
}
