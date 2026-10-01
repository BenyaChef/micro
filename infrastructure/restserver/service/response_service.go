package restserverservice

import (
	"encoding/json"
	"net/http"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

const (
	headerContentType         = "Content-Type"
	headerXContentTypeOptions = "X-Content-Type-Options"
	contentTypeJSON           = "application/json; charset=utf-8"
)

type ResponseService struct {
	errorResponseService *ErrorResponseService
	logPublisher         loggerinterface.LogPublisher
}

func NewResponseService(
	errorResponseService *ErrorResponseService,
	logPublisher loggerinterface.LogPublisher,
) (*ResponseService, error) {
	if errorResponseService == nil {
		return nil, ErrErrorResponseServiceIsRequired
	}

	if logPublisher == nil {
		return nil, ErrLogPublisherIsRequired
	}

	return &ResponseService{errorResponseService: errorResponseService, logPublisher: logPublisher}, nil
}

func (s *ResponseService) JSONResponse(
	writer http.ResponseWriter,
	request *http.Request,
	result any,
	responseCode int,
) {
	body, err := json.Marshal(result)
	if err != nil {
		s.logPublisher.LogError(request.Context(), err)
		s.ErrorResponse(writer, request, restservermodel.ErrWriteResponse)

		return
	}

	s.Response(writer, request, body, responseCode)
}

func (s *ResponseService) Response(
	writer http.ResponseWriter,
	request *http.Request,
	body []byte,
	responseCode int,
) {
	writer.Header().Set(headerContentType, contentTypeJSON)
	writer.Header().Set(headerXContentTypeOptions, "nosniff")
	writer.WriteHeader(responseCode)

	if responseCode == http.StatusNoContent {
		return
	}

	if _, err := writer.Write(body); err != nil {
		s.logPublisher.LogError(request.Context(), err)
	}
}

func (s *ResponseService) ErrorResponse(writer http.ResponseWriter, request *http.Request, err error) {
	s.errorResponseService.ErrorResponse(writer, request, err)
}

func (s *ResponseService) ErrorsResponse(writer http.ResponseWriter, request *http.Request, errs []error) {
	s.errorResponseService.ErrorsResponse(writer, request, errs)
}
