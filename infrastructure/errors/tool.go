package errors

import "errors"

func ContainByCode(err error, errorCode ErrorCode) bool {
	if err == nil {
		return false
	}

	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr.Code() == errorCode
	}

	return false
}

func EqualByCode(err error, errorCode ErrorCode) bool {
	if err == nil {
		return false
	}

	var customErr *Error
	if !errors.As(err, &customErr) {
		return false
	}

	return customErr.Code() == errorCode
}

func CastOrWrap(err error, orErrorCode ErrorCode) *Error {
	if err == nil {
		return nil
	}

	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr
	}

	return NewError(orErrorCode, err.Error())
}

func MessageOf(err error) string {
	if err == nil {
		return ""
	}

	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr.Message()
	}

	return err.Error()
}

func CodeOf(err error) ErrorCode {
	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr.Code()
	}

	return UnknownErrorCode
}
