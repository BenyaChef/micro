package restserverhandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	loggerinterface "micro/infrastructure/logger/interface"
	restservermodel "micro/infrastructure/restserver/model"
)

func RegisterHealth(router chi.Router, serviceName string, logPublisher loggerinterface.LogPublisher) {
	router.Get("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		response := restservermodel.HealthResponse{
			Status:  "ok",
			Service: serviceName,
		}

		if err := json.NewEncoder(writer).Encode(response); err != nil {
			logPublisher.LogError(request.Context(), err)
		}
	})
}
