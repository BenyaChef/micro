package postrestcontroller

import (
	"errors"
	"net/http"

	authgateway "github.com/BenyaChef/micro/blog/adapter/gateway/auth"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	postusecase "github.com/BenyaChef/micro/blog/domain/usecase/post"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
	restserverservice "github.com/BenyaChef/micro/infrastructure/restserver/service"
)

type ErrorResolver struct {
	*restserverservice.DefaultErrorResolver
}

func NewErrorResolver() *ErrorResolver {
	return &ErrorResolver{DefaultErrorResolver: restserverservice.NewDefaultErrorResolver()}
}

func (r *ErrorResolver) GetHTTPCode(err error) int {
	switch {
	case isUnauthorized(err):
		return http.StatusUnauthorized
	case isForbidden(err):
		return http.StatusForbidden
	case isNotFound(err):
		return http.StatusNotFound
	case isBadRequest(err):
		return http.StatusBadRequest
	case isUnavailable(err):
		return http.StatusServiceUnavailable
	default:
		return r.DefaultErrorResolver.GetHTTPCode(err)
	}
}

func isUnauthorized(err error) bool {
	return errors.Is(err, restservermodel.ErrAuthorizationRequired) ||
		errors.Is(err, ErrCallerIsRequired) ||
		apperrors.EqualByCode(err, authgateway.ErrCodeTokenRejected)
}

func isForbidden(err error) bool {
	return apperrors.EqualByCode(err, postusecase.ErrCodeNotAuthor)
}

func isNotFound(err error) bool {
	return apperrors.EqualByCode(err, postusecase.ErrCodePostNotFound)
}

func isBadRequest(err error) bool {
	return errors.Is(err, postentity.ErrTitleIsEmpty) ||
		errors.Is(err, postentity.ErrBodyIsEmpty) ||
		errors.Is(err, postentity.ErrPostIDInvalidFormat) ||
		errors.Is(err, postusecase.ErrInputIsRequired) ||
		errors.Is(err, postusecase.ErrPostIDIsRequired) ||
		apperrors.EqualByCode(err, postentity.ErrCodeTitleTooLong) ||
		apperrors.EqualByCode(err, postentity.ErrCodeBodyTooLong) ||
		apperrors.EqualByCode(err, ErrCodeInvalidQueryParam) ||
		apperrors.EqualByCode(err, restservermodel.ErrCodeUnmarshalRequest) ||
		apperrors.EqualByCode(err, restservermodel.ErrCodeMissingRequiredField)
}

func isUnavailable(err error) bool {
	return apperrors.EqualByCode(err, authgateway.ErrCodeAuthUnavailable)
}
