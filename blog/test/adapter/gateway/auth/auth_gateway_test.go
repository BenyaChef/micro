package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	authgateway "github.com/BenyaChef/micro/blog/adapter/gateway/auth"
	authgatewaymock "github.com/BenyaChef/micro/blog/test/adapter/gateway/auth/mock"
	commonauthgateway "github.com/BenyaChef/micro/common/auth/adapter/gateway"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
)

type AuthGatewayShould struct {
	suite.Suite

	commonGateway *authgatewaymock.AuthGatewayMock
	gateway       *authgateway.AuthGateway
}

func TestAuthGatewayShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AuthGatewayShould))
}

func (s *AuthGatewayShould) SetupTest() {
	s.commonGateway = authgatewaymock.NewAuthGatewayMock()

	gateway, err := authgateway.NewBuilder().
		AuthGateway(s.commonGateway).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.gateway = gateway
}

func (s *AuthGatewayShould) TestBuild_WithoutParams_ReturnError() {
	gateway, err := authgateway.NewBuilder().Build()

	assert.Nil(s.T(), gateway)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, authgateway.ErrAuthGatewayIsRequired))
	assert.True(s.T(), errors.Is(err, authgateway.ErrLogPublisherIsRequired))
}

func (s *AuthGatewayShould) TestValidateToken_AuthRejects_MapsToTokenRejected() {
	s.commonGateway.ValidateTokenFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return commonuserentity.UserID{}, apperrors.NewError("10b754c6-003", "Token is expired")
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "stale")

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, authgateway.ErrCodeTokenRejected))
	assert.True(s.T(), userID.IsZero())
}

func (s *AuthGatewayShould) TestValidateToken_TransportFails_MapsToUnavailable() {
	s.commonGateway.ValidateTokenFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return commonuserentity.UserID{}, commonauthgateway.ErrCallFailed(errors.New("connection refused"))
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "whatever")

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, authgateway.ErrCodeAuthUnavailable))
	assert.False(s.T(), apperrors.EqualByCode(err, authgateway.ErrCodeTokenRejected))
	assert.True(s.T(), userID.IsZero())
}

func (s *AuthGatewayShould) TestValidateToken_Valid_ReturnUserID() {
	expected := commonuserentity.NewUserID()
	s.commonGateway.ValidateTokenFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return expected, nil
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "good")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), expected, userID)
}

func (s *AuthGatewayShould) TestTokenValidator_ReturnsCallerIDAsString() {
	expected := commonuserentity.NewUserID()
	s.commonGateway.ValidateTokenFunc = func(context.Context, string) (commonuserentity.UserID, error) {
		return expected, nil
	}

	validator, err := authgateway.NewTokenValidator(s.gateway)
	require.NoError(s.T(), err)

	callerID, err := validator.ValidateToken(s.T().Context(), "good")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), expected.String(), callerID)
}

func (s *AuthGatewayShould) TestNewTokenValidator_WithoutGateway_ReturnError() {
	validator, err := authgateway.NewTokenValidator(nil)

	assert.Nil(s.T(), validator)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, authgateway.ErrAuthGatewayIsRequired))
}
