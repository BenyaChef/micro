package restserverservice

import (
	"net/http"

	apperrors "micro/infrastructure/errors"
	restservermodel "micro/infrastructure/restserver/model"
)

type DefaultErrorResolver struct{}

func NewDefaultErrorResolver() *DefaultErrorResolver {
	return &DefaultErrorResolver{}
}

func (r *DefaultErrorResolver) GetErrorCode(err error) string {
	return apperrors.CodeOf(err).String()
}

func (r *DefaultErrorResolver) GetErrorText(err error) string {
	return apperrors.MessageOf(err)
}

func (r *DefaultErrorResolver) GetHTTPCode(err error) int {
	switch apperrors.CodeOf(err) {
	case restservermodel.ErrCodeUnmarshalRequest,
		restservermodel.ErrCodeMissingRequiredField:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
