package restserverservice

import (
	"encoding/json"
	"net/http"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

type ErrorResponseService struct {
	errorResolver restserverinterface.ErrorResolver
	logPublisher  loggerinterface.LogPublisher
}

func NewErrorResponseService(
	errorResolver restserverinterface.ErrorResolver,
	logPublisher loggerinterface.LogPublisher,
) (*ErrorResponseService, error) {
	if errorResolver == nil {
		return nil, ErrErrorResolverIsRequired
	}

	if logPublisher == nil {
		return nil, ErrErrorLogPublisherIsRequired
	}

	return &ErrorResponseService{errorResolver: errorResolver, logPublisher: logPublisher}, nil
}

func (s *ErrorResponseService) ErrorResponse(writer http.ResponseWriter, request *http.Request, err error) {
	s.ErrorsResponse(writer, request, []error{err})
}

func (s *ErrorResponseService) ErrorsResponse(writer http.ResponseWriter, request *http.Request, errs []error) {
	if len(errs) == 0 {
		return
	}

	responses := make([]restservermodel.ErrorResponse, 0, len(errs))
	for _, err := range errs {
		responses = append(responses, restservermodel.NewErrorResponse(
			s.errorResolver.GetErrorCode(err),
			s.errorResolver.GetErrorText(err),
		))
	}

	writer.Header().Set(headerContentType, contentTypeJSON)
	writer.Header().Set(headerXContentTypeOptions, "nosniff")
	writer.WriteHeader(s.errorResolver.GetHTTPCode(errs[0]))

	if err := json.NewEncoder(writer).Encode(responses); err != nil {
		s.logPublisher.LogError(request.Context(), err)
	}
}
