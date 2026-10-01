package restservermiddleware

import (
	"net/http"
	"time"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

func LogRequest(logPublisher loggerinterface.LogPublisher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			started := time.Now()
			recorder := newStatusRecorder(writer)

			next.ServeHTTP(recorder, request)

			logPublisher.LogInfo(request.Context(), "request handled",
				"method", request.Method,
				"path", request.URL.Path,
				"status", recorder.Status(),
				"took", time.Since(started),
				"request_id", restservermodel.RequestIDFrom(request.Context()),
			)
		})
	}
}
