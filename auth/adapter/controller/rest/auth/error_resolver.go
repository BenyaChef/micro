package authrestcontroller

import (
	"errors"
	"net/http"

	userrepository "github.com/BenyaChef/micro/auth/adapter/repository/user"
	tokenusecase "github.com/BenyaChef/micro/auth/domain/usecase/token"
	userusecase "github.com/BenyaChef/micro/auth/domain/usecase/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	"github.com/BenyaChef/micro/infrastructure/jwt"
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
	case isConflict(err):
		return http.StatusConflict
	case isNotFound(err):
		return http.StatusNotFound
	case isBadRequest(err):
		return http.StatusBadRequest
	default:
		return r.DefaultErrorResolver.GetHTTPCode(err)
	}
}

func isUnauthorized(err error) bool {
	return errors.Is(err, ErrAuthorizationRequired) ||
		errors.Is(err, userusecase.ErrInvalidCredentials) ||
		errors.Is(err, tokenusecase.ErrTokenIsRequired) ||
		errors.Is(err, jwt.ErrTokenInvalid) ||
		errors.Is(err, jwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrSubjectIsEmpty) ||
		errors.Is(err, commonuserentity.ErrUserIDInvalidFormat)
}

func isConflict(err error) bool {
	return apperrors.EqualByCode(err, userusecase.ErrCodeEmailAlreadyTaken) ||
		apperrors.EqualByCode(err, userrepository.ErrCodeEmailConflict)
}

func isNotFound(err error) bool {
	return apperrors.EqualByCode(err, userusecase.ErrCodeUserNotFound)
}

func isBadRequest(err error) bool {
	return errors.Is(err, emailprimitive.ErrEmailIsEmpty) ||
		errors.Is(err, emailprimitive.ErrEmailInvalidFormat) ||
		errors.Is(err, userusecase.ErrDataIsRequired) ||
		errors.Is(err, userusecase.ErrPasswordIsRequired) ||
		apperrors.EqualByCode(err, userusecase.ErrCodePasswordTooShort)
}
