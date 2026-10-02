package jwt_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	"github.com/BenyaChef/micro/infrastructure/jwt"
)

const (
	testSecret      = "0123456789abcdef0123456789abcdef"
	testOtherSecret = "fedcba9876543210fedcba9876543210"
	testIssuer      = "auth"
	testSubject     = "550e8400-e29b-41d4-a716-446655440000"
)

type ServiceShould struct {
	suite.Suite
}

func TestServiceShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceShould))
}

func (s *ServiceShould) newService(secret string, ttl time.Duration) *jwt.Service {
	service, err := jwt.NewBuilder().Secret(secret).Issuer(testIssuer).TTL(ttl).Build()
	require.NoError(s.T(), err)

	return service
}

func (s *ServiceShould) TestBuild_WithoutParams_ReturnError() {
	service, err := jwt.NewBuilder().Build()

	assert.Nil(s.T(), service)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, jwt.ErrSecretIsRequired))
	assert.True(s.T(), errors.Is(err, jwt.ErrIssuerIsRequired))
}

func (s *ServiceShould) TestBuild_ShortSecret_ReturnError() {
	service, err := jwt.NewBuilder().Secret("too-short").Issuer(testIssuer).Build()

	assert.Nil(s.T(), service)
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, jwt.ErrCodeSecretTooShort))
}

func (s *ServiceShould) TestIssue_EmptySubject_ReturnError() {
	service := s.newService(testSecret, time.Hour)

	rawToken, expiresIn, err := service.Issue("")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, jwt.ErrSubjectIsEmpty))
	assert.Empty(s.T(), rawToken)
	assert.Zero(s.T(), expiresIn)
}

func (s *ServiceShould) TestIssueThenSubject_ReturnSameSubject() {
	service := s.newService(testSecret, time.Hour)

	rawToken, expiresIn, err := service.Issue(testSubject)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), time.Hour, expiresIn)

	subject, err := service.Subject(rawToken)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), testSubject, subject)
}

func (s *ServiceShould) TestSubject_MalformedToken_ReturnError() {
	service := s.newService(testSecret, time.Hour)

	subject, err := service.Subject("not.a.token")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, jwt.ErrTokenInvalid))
	assert.Empty(s.T(), subject)
}

func (s *ServiceShould) TestSubject_ForeignSecret_ReturnError() {
	issuer := s.newService(testSecret, time.Hour)
	verifier := s.newService(testOtherSecret, time.Hour)

	rawToken, _, err := issuer.Issue(testSubject)
	require.NoError(s.T(), err)

	subject, err := verifier.Subject(rawToken)

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, jwt.ErrTokenInvalid))
	assert.Empty(s.T(), subject)
}

func (s *ServiceShould) TestSubject_ExpiredToken_ReturnExpiredError() {
	service := s.newService(testSecret, -time.Second)

	rawToken, _, err := service.Issue(testSubject)
	require.NoError(s.T(), err)

	subject, err := service.Subject(rawToken)

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, jwt.ErrTokenExpired))
	assert.Empty(s.T(), subject)
}

func (s *ServiceShould) TestSubject_ForeignIssuer_ReturnError() {
	issuer, err := jwt.NewBuilder().Secret(testSecret).Issuer("blog").TTL(time.Hour).Build()
	require.NoError(s.T(), err)

	verifier := s.newService(testSecret, time.Hour)

	rawToken, _, err := issuer.Issue(testSubject)
	require.NoError(s.T(), err)

	subject, err := verifier.Subject(rawToken)

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, jwt.ErrTokenInvalid))
	assert.Empty(s.T(), subject)
}
