package restservermiddleware

import (
	"net/http"
	"time"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

func LogRequest(logPublisher loggerinterface.LogPublisher, excludedPaths ...string) func(http.Handler) http.Handler {
	excluded := make(map[string]bool, len(excludedPaths))
	for _, path := range excludedPaths {
		excluded[path] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if excluded[request.URL.Path] {
				next.ServeHTTP(writer, request)

				return
			}

			started := time.Now()
			recorder := newStatusRecorder(writer)

			next.ServeHTTP(recorder, request)

			attrs := []any{
				"method", request.Method,
				"path", request.URL.Path,
				"status", recorder.Status(),
				"took", time.Since(started).String(),
				"request_id", restservermodel.RequestIDFrom(request.Context()),
			}

			logByStatus(request, logPublisher, recorder.Status(), attrs)
		})
	}
}

func logByStatus(
	request *http.Request,
	logPublisher loggerinterface.LogPublisher,
	status int,
	attrs []any,
) {
	const message = "request handled"

	switch {
	case status >= http.StatusInternalServerError:
		logPublisher.LogErrorMsg(request.Context(), message, attrs...)
	case status >= http.StatusBadRequest:
		logPublisher.LogWarn(request.Context(), message, attrs...)
	default:
		logPublisher.LogInfo(request.Context(), message, attrs...)
	}
}
