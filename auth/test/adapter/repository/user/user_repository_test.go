package user_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	userrepository "github.com/BenyaChef/micro/auth/adapter/repository/user"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	userentitystub "github.com/BenyaChef/micro/auth/test/domain/entity/user/stub"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	"github.com/BenyaChef/micro/infrastructure/postgres"
)

const envTestDSN = "POSTGRES_TEST_DSN"

type UserRepositoryShould struct {
	suite.Suite

	client     *postgres.Client
	repository *userrepository.UserRepository
}

func TestUserRepositoryShould(t *testing.T) {
	suite.Run(t, new(UserRepositoryShould))
}

func (s *UserRepositoryShould) SetupSuite() {
	dsn := os.Getenv(envTestDSN)
	if dsn == "" {
		s.T().Skipf("%s is not set, skipping test against live database", envTestDSN)
	}

	client, err := postgres.NewBuilder().DSN(dsn).Build(s.T().Context())
	require.NoError(s.T(), err)

	s.client = client

	repository, err := userrepository.NewBuilder().
		Executor(client).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.repository = repository
}

func (s *UserRepositoryShould) TearDownSuite() {
	if s.client != nil {
		s.client.Close()
	}
}

func (s *UserRepositoryShould) insertFreshUser() *userentity.User {
	s.T().Helper()

	email, err := emailprimitive.EmailFrom(fmt.Sprintf("user-%s@example.com", commonuserentity.NewUserID()))
	require.NoError(s.T(), err)

	user, err := userentity.NewBuilder().
		ID(commonuserentity.NewUserID()).
		Email(email).
		PasswordHash(userentitystub.GetPasswordHash()).
		Build()
	require.NoError(s.T(), err)

	require.NoError(s.T(), s.repository.Insert(s.T().Context(), user))

	s.T().Cleanup(func() {
		cleanupCtx := context.Background()
		_, cleanupErr := s.client.Exec(cleanupCtx, "DELETE FROM users WHERE id = $1", user.ID().String())
		require.NoError(s.T(), cleanupErr)
	})

	return user
}

func (s *UserRepositoryShould) TestBuild_WithoutParams_ReturnError() {
	repository, err := userrepository.NewBuilder().Build()

	assert.Nil(s.T(), repository)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userrepository.ErrExecutorIsRequired))
	assert.True(s.T(), errors.Is(err, userrepository.ErrLogPublisherIsRequired))
}

func (s *UserRepositoryShould) TestInsert_NilUser_ReturnError() {
	err := s.repository.Insert(s.T().Context(), nil)

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userrepository.ErrUserIsRequired))
}

func (s *UserRepositoryShould) TestGetByID_InsertedUser_ReturnEqualEntity() {
	inserted := s.insertFreshUser()

	found, err := s.repository.GetByID(s.T().Context(), inserted.ID())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), inserted.ID(), found.ID())
	assert.Equal(s.T(), inserted.Email(), found.Email())
	assert.Equal(s.T(), inserted.PasswordHash(), found.PasswordHash())
	assert.Equal(s.T(), inserted.CreatedAt(), found.CreatedAt())
	assert.Equal(s.T(), inserted.UpdatedAt(), found.UpdatedAt())
}

func (s *UserRepositoryShould) TestGetByEmail_IgnoresCase() {
	inserted := s.insertFreshUser()

	upperEmail, err := emailprimitive.EmailFrom(strings.ToUpper(inserted.Email().String()))
	require.NoError(s.T(), err)

	found, err := s.repository.GetByEmail(s.T().Context(), upperEmail)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), inserted.ID(), found.ID())
}

func (s *UserRepositoryShould) TestGetByID_UnknownUser_ReturnNilNil() {
	found, err := s.repository.GetByID(s.T().Context(), commonuserentity.NewUserID())

	require.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserRepositoryShould) TestGetByEmail_UnknownUser_ReturnNilNil() {
	email, err := emailprimitive.EmailFrom("nobody@example.com")
	require.NoError(s.T(), err)

	found, err := s.repository.GetByEmail(s.T().Context(), email)

	require.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UserRepositoryShould) TestInsert_DuplicateEmail_ReturnConflict() {
	inserted := s.insertFreshUser()

	duplicate, err := userentity.NewBuilder().
		ID(commonuserentity.NewUserID()).
		Email(inserted.Email()).
		PasswordHash(userentitystub.GetPasswordHash()).
		Build()
	require.NoError(s.T(), err)

	err = s.repository.Insert(s.T().Context(), duplicate)

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, userrepository.ErrCodeEmailConflict))
}
