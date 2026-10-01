package authrestcontroller

import (
	"net/http"
	"strings"

	requestrestmodel "github.com/BenyaChef/micro/auth/adapter/controller/rest/auth/request"
	responserestmodel "github.com/BenyaChef/micro/auth/adapter/controller/rest/auth/response"
	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
)

const bearerPrefix = "Bearer "

type AuthController struct {
	*restservercontroller.BaseController

	userUseCase  usecaseinterface.UserUseCase
	tokenUseCase usecaseinterface.TokenUseCase
}

func (c *AuthController) Register(writer http.ResponseWriter, request *http.Request) {
	requestModel := &requestrestmodel.RegisterRequest{}
	if err := c.FillAndValidateReqModel(request, requestModel); err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	data := &boundarymodel.RegisterData{Email: requestModel.Email, Password: requestModel.Password}

	user, err := c.userUseCase.Register(request.Context(), data)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, responserestmodel.FromUser(user), http.StatusCreated)
}

func (c *AuthController) Login(writer http.ResponseWriter, request *http.Request) {
	requestModel := &requestrestmodel.LoginRequest{}
	if err := c.FillAndValidateReqModel(request, requestModel); err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	data := &boundarymodel.LoginData{Email: requestModel.Email, Password: requestModel.Password}

	result, err := c.userUseCase.Login(request.Context(), data)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, responserestmodel.FromLoginResult(result), http.StatusOK)
}

func (c *AuthController) Me(writer http.ResponseWriter, request *http.Request) {
	rawToken, err := bearerToken(request)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	userID, err := c.tokenUseCase.Validate(request.Context(), rawToken)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	user, err := c.userUseCase.GetByID(request.Context(), userID)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, responserestmodel.FromUser(user), http.StatusOK)
}

func bearerToken(request *http.Request) (string, error) {
	header := request.Header.Get("Authorization")
	if !strings.HasPrefix(header, bearerPrefix) {
		return "", ErrAuthorizationRequired
	}

	rawToken := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
	if rawToken == "" {
		return "", ErrAuthorizationRequired
	}

	return rawToken, nil
}
