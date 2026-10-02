package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	authrestcontroller "github.com/BenyaChef/micro/auth/adapter/controller/rest/auth"
	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	userusecase "github.com/BenyaChef/micro/auth/domain/usecase/user"
	authrestcontrollermock "github.com/BenyaChef/micro/auth/test/adapter/controller/rest/auth/mock"
	userentitystub "github.com/BenyaChef/micro/auth/test/domain/entity/user/stub"
	userusecasemock "github.com/BenyaChef/micro/auth/test/domain/usecase/user/mock"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	"github.com/BenyaChef/micro/infrastructure/jwt"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
	restserverservice "github.com/BenyaChef/micro/infrastructure/restserver/service"
)

type AuthControllerShould struct {
	suite.Suite

	userUseCase  *authrestcontrollermock.UserUseCaseMock
	tokenUseCase *userusecasemock.TokenUseCaseMock
	controller   *authrestcontroller.AuthController
}

func TestAuthControllerShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AuthControllerShould))
}

func (s *AuthControllerShould) SetupTest() {
	logPublisher := loggerteststub.NewLogPublisherStub()

	errorResponseService, err := restserverservice.NewErrorResponseService(
		authrestcontroller.NewErrorResolver(), logPublisher,
	)
	require.NoError(s.T(), err)

	responseService, err := restserverservice.NewResponseService(errorResponseService, logPublisher)
	require.NoError(s.T(), err)

	baseController, err := restservercontroller.NewBaseController(responseService, logPublisher)
	require.NoError(s.T(), err)

	s.userUseCase = authrestcontrollermock.NewUserUseCaseMock()
	s.tokenUseCase = userusecasemock.NewTokenUseCaseMock()

	controller, err := authrestcontroller.NewBuilder().
		BaseController(baseController).
		UserUseCase(s.userUseCase).
		TokenUseCase(s.tokenUseCase).
		Build()
	require.NoError(s.T(), err)

	s.controller = controller
}

func (s *AuthControllerShould) call(
	handler http.HandlerFunc,
	method, target, body string,
	header map[string]string,
) *httptest.ResponseRecorder {
	s.T().Helper()

	request := httptest.NewRequestWithContext(s.T().Context(), method, target, strings.NewReader(body))
	for key, value := range header {
		request.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()
	handler(recorder, request)

	return recorder
}

func (s *AuthControllerShould) TestBuild_WithoutParams_ReturnError() {
	controller, err := authrestcontroller.NewBuilder().Build()

	assert.Nil(s.T(), controller)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, authrestcontroller.ErrBaseControllerIsRequired))
	assert.True(s.T(), errors.Is(err, authrestcontroller.ErrUserUseCaseIsRequired))
	assert.True(s.T(), errors.Is(err, authrestcontroller.ErrTokenUseCaseIsRequired))
}

func (s *AuthControllerShould) TestRegister_ValidBody_Return201WithUser() {
	expected := userentitystub.GetValidUser()
	s.userUseCase.RegisterFunc = func(context.Context, *boundarymodel.RegisterData) (*userentity.User, error) {
		return expected, nil
	}

	recorder := s.call(s.controller.Register, http.MethodPost, "/auth/register",
		`{"email":"user@example.com","password":"correct horse battery staple"}`, nil)

	require.Equal(s.T(), http.StatusCreated, recorder.Code)

	var body map[string]any
	require.NoError(s.T(), json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(s.T(), expected.ID().String(), body["user_id"])
	assert.Equal(s.T(), expected.Email().String(), body["email"])
}

func (s *AuthControllerShould) TestRegister_MalformedJSON_Return400() {
	recorder := s.call(s.controller.Register, http.MethodPost, "/auth/register", `{oops`, nil)

	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

func (s *AuthControllerShould) TestRegister_MissingPassword_Return400() {
	recorder := s.call(s.controller.Register, http.MethodPost, "/auth/register",
		`{"email":"user@example.com"}`, nil)

	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

func (s *AuthControllerShould) TestRegister_EmailTaken_Return409() {
	s.userUseCase.RegisterFunc = func(context.Context, *boundarymodel.RegisterData) (*userentity.User, error) {
		return nil, userusecase.ErrEmailAlreadyTaken("user@example.com")
	}

	recorder := s.call(s.controller.Register, http.MethodPost, "/auth/register",
		`{"email":"user@example.com","password":"correct horse battery staple"}`, nil)

	assert.Equal(s.T(), http.StatusConflict, recorder.Code)
}

func (s *AuthControllerShould) TestLogin_ValidCredentials_Return200WithToken() {
	s.userUseCase.LoginFunc = func(context.Context, *boundarymodel.LoginData) (*boundarymodel.LoginResult, error) {
		return &boundarymodel.LoginResult{AccessToken: "token-value", ExpiresIn: 3600}, nil
	}

	recorder := s.call(s.controller.Login, http.MethodPost, "/auth/login",
		`{"email":"user@example.com","password":"correct horse battery staple"}`, nil)

	require.Equal(s.T(), http.StatusOK, recorder.Code)

	var body map[string]any
	require.NoError(s.T(), json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(s.T(), "token-value", body["access_token"])
	assert.InDelta(s.T(), 3600, body["expires_in"], 0)
}

func (s *AuthControllerShould) TestLogin_InvalidCredentials_Return401() {
	s.userUseCase.LoginFunc = func(context.Context, *boundarymodel.LoginData) (*boundarymodel.LoginResult, error) {
		return nil, userusecase.ErrInvalidCredentials
	}

	recorder := s.call(s.controller.Login, http.MethodPost, "/auth/login",
		`{"email":"user@example.com","password":"nope"}`, nil)

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *AuthControllerShould) TestMe_WithoutAuthorizationHeader_Return401() {
	recorder := s.call(s.controller.Me, http.MethodGet, "/auth/me", "", nil)

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *AuthControllerShould) TestMe_NonBearerHeader_Return401() {
	recorder := s.call(s.controller.Me, http.MethodGet, "/auth/me", "",
		map[string]string{"Authorization": "Basic dXNlcjpwYXNz"})

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *AuthControllerShould) TestMe_ExpiredToken_Return401() {
	s.tokenUseCase.ValidateFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return commonuserentity.UserID{}, jwt.ErrTokenExpired
	}

	recorder := s.call(s.controller.Me, http.MethodGet, "/auth/me", "",
		map[string]string{"Authorization": "Bearer stale"})

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *AuthControllerShould) TestMe_EmptyBearerValue_Return401() {
	recorder := s.call(s.controller.Me, http.MethodGet, "/auth/me", "",
		map[string]string{"Authorization": "Bearer   "})

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *AuthControllerShould) TestMe_UserNotFound_Return404() {
	s.userUseCase.GetByIDFunc = func(context.Context, commonuserentity.UserID) (*userentity.User, error) {
		return nil, userusecase.ErrUserNotFound("some-id")
	}

	recorder := s.call(s.controller.Me, http.MethodGet, "/auth/me", "",
		map[string]string{"Authorization": "Bearer good"})

	assert.Equal(s.T(), http.StatusNotFound, recorder.Code)
}

func (s *AuthControllerShould) TestMe_ValidToken_Return200WithUser() {
	expected := userentitystub.GetValidUser()
	s.userUseCase.GetByIDFunc = func(context.Context, commonuserentity.UserID) (*userentity.User, error) {
		return expected, nil
	}

	recorder := s.call(s.controller.Me, http.MethodGet, "/auth/me", "",
		map[string]string{"Authorization": "Bearer good"})

	require.Equal(s.T(), http.StatusOK, recorder.Code)

	var body map[string]any
	require.NoError(s.T(), json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(s.T(), expected.ID().String(), body["user_id"])
}
