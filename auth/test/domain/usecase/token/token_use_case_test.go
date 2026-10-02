package token_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	tokenusecase "github.com/BenyaChef/micro/auth/domain/usecase/token"
	tokenusecasemock "github.com/BenyaChef/micro/auth/test/domain/usecase/token/mock"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
)

type TokenUseCaseShould struct {
	suite.Suite

	tokenService *tokenusecasemock.TokenServiceMock
	useCase      *tokenusecase.TokenUseCase
}

func TestTokenUseCaseShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TokenUseCaseShould))
}

func (s *TokenUseCaseShould) SetupTest() {
	s.tokenService = tokenusecasemock.NewTokenServiceMock()

	useCase, err := tokenusecase.NewBuilder().
		TokenService(s.tokenService).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.useCase = useCase
}

func (s *TokenUseCaseShould) TestBuild_WithoutParams_ReturnError() {
	useCase, err := tokenusecase.NewBuilder().Build()

	assert.Nil(s.T(), useCase)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, tokenusecase.ErrTokenServiceIsRequired))
	assert.True(s.T(), errors.Is(err, tokenusecase.ErrLogPublisherIsRequired))
}

func (s *TokenUseCaseShould) TestIssue_ZeroUserID_ReturnError() {
	result, err := s.useCase.Issue(s.T().Context(), commonuserentity.UserID{})

	assert.Nil(s.T(), result)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, tokenusecase.ErrUserIDIsRequired))
}

func (s *TokenUseCaseShould) TestIssue_ValidUserID_ReturnTokenWithTTLInSeconds() {
	result, err := s.useCase.Issue(s.T().Context(), commonuserentity.NewUserID())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.Equal(s.T(), tokenusecasemock.StubRawToken, result.AccessToken)
	assert.Equal(s.T(), int(tokenusecasemock.StubExpiresIn.Seconds()), result.ExpiresIn)
}

func (s *TokenUseCaseShould) TestIssue_ServiceFails_WrapError() {
	cause := errors.New("signing broken")
	s.tokenService.IssueFunc = func(string) (string, time.Duration, error) { return "", 0, cause }

	result, err := s.useCase.Issue(s.T().Context(), commonuserentity.NewUserID())

	assert.Nil(s.T(), result)
	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, cause)
}

func (s *TokenUseCaseShould) TestValidate_EmptyToken_ReturnError() {
	userID, err := s.useCase.Validate(s.T().Context(), "")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, tokenusecase.ErrTokenIsRequired))
	assert.True(s.T(), userID.IsZero())
}

func (s *TokenUseCaseShould) TestValidate_ServiceRejects_PassErrorThrough() {
	cause := errors.New("token is rotten")
	s.tokenService.SubjectFunc = func(string) (string, error) { return "", cause }

	userID, err := s.useCase.Validate(s.T().Context(), "whatever")

	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, cause)
	assert.True(s.T(), userID.IsZero())
}

func (s *TokenUseCaseShould) TestValidate_SubjectIsNotUUID_ReturnError() {
	s.tokenService.SubjectFunc = func(string) (string, error) { return "not-a-uuid", nil }

	userID, err := s.useCase.Validate(s.T().Context(), "whatever")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, commonuserentity.ErrUserIDInvalidFormat))
	assert.True(s.T(), userID.IsZero())
}

func (s *TokenUseCaseShould) TestValidate_ValidToken_ReturnUserID() {
	userID, err := s.useCase.Validate(s.T().Context(), "whatever")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), "550e8400-e29b-41d4-a716-446655440000", userID.String())
}
