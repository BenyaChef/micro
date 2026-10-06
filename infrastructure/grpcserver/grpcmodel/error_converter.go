package grpcmodel

import "github.com/BenyaChef/micro/infrastructure/errors"

func ToErrorModel(err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		Code: errors.CodeOf(err).String(),
		Msg:  errors.MessageOf(err),
	}
}

func FromErrorModel(model *Error) error {
	if model == nil {
		return nil
	}

	return errors.NewError(errors.ErrorCode(model.GetCode()), model.GetMsg())
}
