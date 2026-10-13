package restservermiddleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	infraotel "github.com/BenyaChef/micro/infrastructure/otel"
)

const unmatchedRoute = "unmatched"

func Metrics(monitor *infraotel.HTTPMonitor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			started := time.Now()
			statusWriter := newStatusRecorder(writer)

			next.ServeHTTP(statusWriter, request)

			monitor.RecordRequest(
				request.Context(),
				request.Method,
				routePattern(request),
				statusWriter.Status(),
				time.Since(started).Seconds(),
			)
		})
	}
}

func routePattern(request *http.Request) string {
	routeContext := chi.RouteContext(request.Context())
	if routeContext == nil {
		return unmatchedRoute
	}

	pattern := routeContext.RoutePattern()
	if pattern == "" {
		return unmatchedRoute
	}

	return pattern
}
