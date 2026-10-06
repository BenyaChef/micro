package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	authgrpccontroller "github.com/BenyaChef/micro/auth/adapter/controller/grpc/auth"
	tokenusecase "github.com/BenyaChef/micro/auth/domain/usecase/token"
	userusecasemock "github.com/BenyaChef/micro/auth/test/domain/usecase/user/mock"
	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	"github.com/BenyaChef/micro/infrastructure/jwt"
)

type AuthGRPCControllerShould struct {
	suite.Suite

	tokenUseCase *userusecasemock.TokenUseCaseMock
	controller   *authgrpccontroller.AuthController
}

func TestAuthGRPCControllerShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AuthGRPCControllerShould))
}

func (s *AuthGRPCControllerShould) SetupTest() {
	s.tokenUseCase = userusecasemock.NewTokenUseCaseMock()

	controller, err := authgrpccontroller.NewBuilder().
		TokenUseCase(s.tokenUseCase).
		Build()
	require.NoError(s.T(), err)

	s.controller = controller
}

func (s *AuthGRPCControllerShould) validate(rawToken string) *authproto.ValidateTokenResponse {
	s.T().Helper()

	response, err := s.controller.ValidateToken(
		s.T().Context(),
		&authproto.ValidateTokenRequest{AccessToken: rawToken},
	)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), response)

	return response
}

func (s *AuthGRPCControllerShould) TestBuild_WithoutParams_ReturnError() {
	controller, err := authgrpccontroller.NewBuilder().Build()

	assert.Nil(s.T(), controller)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, authgrpccontroller.ErrTokenUseCaseIsRequired))
}

func (s *AuthGRPCControllerShould) TestValidateToken_ValidToken_ReturnUserID() {
	expected := commonuserentity.NewUserID()
	s.tokenUseCase.ValidateFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return expected, nil
	}

	response := s.validate("good-token")

	require.Nil(s.T(), response.GetError())
	require.NotNil(s.T(), response.GetToken())
	assert.Equal(s.T(), expected.String(), response.GetToken().GetUserId())
}

func (s *AuthGRPCControllerShould) TestValidateToken_ExpiredToken_ReturnErrorInsideResponse() {
	s.tokenUseCase.ValidateFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return commonuserentity.UserID{}, jwt.ErrTokenExpired
	}

	response := s.validate("stale")

	require.Nil(s.T(), response.GetToken())
	require.NotNil(s.T(), response.GetError())
	assert.Equal(s.T(), string(jwt.ErrCodeTokenExpired), response.GetError().GetCode())
}

func (s *AuthGRPCControllerShould) TestValidateToken_EmptyToken_ReturnErrorInsideResponse() {
	s.tokenUseCase.ValidateFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return commonuserentity.UserID{}, tokenusecase.ErrTokenIsRequired
	}

	response := s.validate("")

	require.Nil(s.T(), response.GetToken())
	require.NotNil(s.T(), response.GetError())
	assert.NotEmpty(s.T(), response.GetError().GetMsg())
}
