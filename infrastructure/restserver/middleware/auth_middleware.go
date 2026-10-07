package restservermiddleware

import (
	"net/http"
	"strings"

	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
	restserverservice "github.com/BenyaChef/micro/infrastructure/restserver/service"
)

const bearerPrefix = "Bearer "

func Auth(
	validator restserverinterface.TokenValidator,
	responseService *restserverservice.ResponseService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			rawToken := bearerToken(request)
			if rawToken == "" {
				responseService.ErrorResponse(writer, request, restservermodel.ErrAuthorizationRequired)

				return
			}

			callerID, err := validator.ValidateToken(request.Context(), rawToken)
			if err != nil {
				responseService.ErrorResponse(writer, request, err)

				return
			}

			ctx := restservermodel.WithCallerID(request.Context(), callerID)

			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func bearerToken(request *http.Request) string {
	header := request.Header.Get("Authorization")
	if !strings.HasPrefix(header, bearerPrefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
}
