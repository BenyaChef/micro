package restservercontroller

import (
	"errors"
	"io"
	"net/http"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
	restserverservice "github.com/BenyaChef/micro/infrastructure/restserver/service"
)

type BaseController struct {
	responseService *restserverservice.ResponseService
	logPublisher    loggerinterface.LogPublisher
}

func NewBaseController(
	responseService *restserverservice.ResponseService,
	logPublisher loggerinterface.LogPublisher,
) (*BaseController, error) {
	var errs []error

	if responseService == nil {
		errs = append(errs, ErrResponseServiceIsRequired)
	}

	if logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &BaseController{responseService: responseService, logPublisher: logPublisher}, nil
}

func (c *BaseController) FillAndValidateReqModel(request *http.Request, requestModel restserverinterface.ValidateRequestModel) error {
	if err := c.FillReqModel(request, requestModel); err != nil {
		return err
	}

	return requestModel.ValidateRequest()
}

func (c *BaseController) FillReqModel(request *http.Request, requestModel restserverinterface.RequestModel) error {
	body, err := c.GetReqBody(request)
	if err != nil {
		return err
	}

	if err := requestModel.FillFromBytes(body); err != nil {
		return restservermodel.ErrUnmarshalRequest(err.Error())
	}

	return nil
}

func (c *BaseController) GetReqBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, restservermodel.ErrUnmarshalRequest("request body is nil")
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, restservermodel.ErrUnmarshalRequest(err.Error())
	}

	return body, nil
}

func (c *BaseController) JSONResponse(writer http.ResponseWriter, request *http.Request, result any, responseCode int) {
	c.responseService.JSONResponse(writer, request, result, responseCode)
}

func (c *BaseController) ErrorResponse(writer http.ResponseWriter, request *http.Request, err error) {
	c.responseService.ErrorResponse(writer, request, err)
}

func (c *BaseController) ErrorResponseWithLog(writer http.ResponseWriter, request *http.Request, errs ...error) {
	c.logPublisher.LogError(request.Context(), errs...)
	c.responseService.ErrorsResponse(writer, request, errs)
}
