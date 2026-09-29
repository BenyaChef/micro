package restservermiddleware

import (
	"net/http"

	"github.com/google/uuid"

	restservermodel "micro/infrastructure/restserver/model"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestID := request.Header.Get(restservermodel.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		writer.Header().Set(restservermodel.RequestIDHeader, requestID)

		next.ServeHTTP(writer, request.WithContext(
			restservermodel.WithRequestID(request.Context(), requestID),
		))
	})
}
