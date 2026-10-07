package gateway_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	commonauthgateway "github.com/BenyaChef/micro/common/auth/adapter/gateway"
	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	commonauthgatewaymock "github.com/BenyaChef/micro/common/auth/test/adapter/gateway/mock"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	"github.com/BenyaChef/micro/infrastructure/grpcserver/grpcmodel"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
)

type AuthGRPCGatewayShould struct {
	suite.Suite

	client  *commonauthgatewaymock.AuthServiceClientMock
	gateway *commonauthgateway.AuthGRPCGateway
}

func TestAuthGRPCGatewayShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AuthGRPCGatewayShould))
}

func (s *AuthGRPCGatewayShould) SetupTest() {
	s.client = commonauthgatewaymock.NewAuthServiceClientMock()

	gateway, err := commonauthgateway.NewBuilder().
		Client(s.client).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.gateway = gateway
}

func (s *AuthGRPCGatewayShould) TestBuild_WithoutParams_ReturnError() {
	gateway, err := commonauthgateway.NewBuilder().Build()

	assert.Nil(s.T(), gateway)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, commonauthgateway.ErrClientIsRequired))
	assert.True(s.T(), errors.Is(err, commonauthgateway.ErrLogPublisherIsRequired))
}

func (s *AuthGRPCGatewayShould) TestValidateToken_EmptyToken_DoesNotCallService() {
	called := false
	s.client.ValidateTokenFunc = func(context.Context, *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error) {
		called = true

		return nil, nil
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, commonauthgateway.ErrTokenIsRequired))
	assert.True(s.T(), userID.IsZero())
	assert.False(s.T(), called, "пустой токен не должен уходить по сети")
}

func (s *AuthGRPCGatewayShould) TestValidateToken_TransportFails_WrapError() {
	cause := errors.New("connection refused")
	s.client.ValidateTokenFunc = func(context.Context, *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error) {
		return nil, cause
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "whatever")

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, cause)
	assert.True(s.T(), apperrors.EqualByCode(err, commonauthgateway.ErrCodeCallFailed))
	assert.True(s.T(), userID.IsZero())
}

func (s *AuthGRPCGatewayShould) TestValidateToken_ServiceReturnsError_PreserveCode() {
	s.client.ValidateTokenFunc = func(context.Context, *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error) {
		return &authproto.ValidateTokenResponse{
			Result: &authproto.ValidateTokenResponse_Error{
				Error: &grpcmodel.Error{Code: "10b754c6-003", Msg: "Token is expired"},
			},
		}, nil
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "stale")

	require.Error(s.T(), err)
	assert.Equal(s.T(), apperrors.ErrorCode("10b754c6-003"), apperrors.CodeOf(err))
	assert.True(s.T(), userID.IsZero())
}

func (s *AuthGRPCGatewayShould) TestValidateToken_EmptyResult_ReturnError() {
	s.client.ValidateTokenFunc = func(context.Context, *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error) {
		return &authproto.ValidateTokenResponse{}, nil
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "whatever")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, commonauthgateway.ErrEmptyResult))
	assert.True(s.T(), userID.IsZero())
}

func (s *AuthGRPCGatewayShould) TestValidateToken_MalformedUserID_ReturnError() {
	s.client.ValidateTokenFunc = func(context.Context, *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error) {
		return &authproto.ValidateTokenResponse{
			Result: &authproto.ValidateTokenResponse_Token{
				Token: &authproto.ValidatedToken{UserId: "not-a-uuid"},
			},
		}, nil
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "whatever")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, commonuserentity.ErrUserIDInvalidFormat))
	assert.True(s.T(), userID.IsZero())
}

func (s *AuthGRPCGatewayShould) TestValidateToken_ValidResponse_ReturnUserID() {
	expected := commonuserentity.NewUserID()
	s.client.ValidateTokenFunc = func(context.Context, *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error) {
		return &authproto.ValidateTokenResponse{
			Result: &authproto.ValidateTokenResponse_Token{
				Token: &authproto.ValidatedToken{UserId: expected.String()},
			},
		}, nil
	}

	userID, err := s.gateway.ValidateToken(s.T().Context(), "good")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), expected, userID)
}
