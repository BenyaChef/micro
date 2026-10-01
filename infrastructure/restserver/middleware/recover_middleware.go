package restservermiddleware

import (
	"net/http"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

func Recover(logPublisher loggerinterface.LogPublisher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}

				logPublisher.LogWarn(request.Context(), "panic recovered",
					"panic", recovered,
					"method", request.Method,
					"path", request.URL.Path,
					"request_id", restservermodel.RequestIDFrom(request.Context()),
				)

				writer.WriteHeader(http.StatusInternalServerError)
			}()

			next.ServeHTTP(writer, request)
		})
	}
}
