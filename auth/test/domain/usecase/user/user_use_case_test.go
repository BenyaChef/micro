package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	userusecase "github.com/BenyaChef/micro/auth/domain/usecase/user"
	userentitystub "github.com/BenyaChef/micro/auth/test/domain/entity/user/stub"
	userusecasemock "github.com/BenyaChef/micro/auth/test/domain/usecase/user/mock"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
)

const validPassword = "correct horse battery staple"

type UserUseCaseShould struct {
	suite.Suite

	userRepository *userusecasemock.UserRepositoryMock
	passwordHasher *userusecasemock.PasswordHasherMock
	tokenUseCase   *userusecasemock.TokenUseCaseMock
	useCase        *userusecase.UserUseCase
}

func TestUserUseCaseShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserUseCaseShould))
}

func (s *UserUseCaseShould) SetupTest() {
	s.userRepository = userusecasemock.NewUserRepositoryMock()
	s.passwordHasher = userusecasemock.NewPasswordHasherMock()
	s.tokenUseCase = userusecasemock.NewTokenUseCaseMock()

	useCase, err := userusecase.NewBuilder().
		UserRepository(s.userRepository).
		PasswordHasher(s.passwordHasher).
		TokenUseCase(s.tokenUseCase).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.useCase = useCase
}

func (s *UserUseCaseShould) registerData() *boundarymodel.RegisterData {
	return &boundarymodel.RegisterData{Email: userentitystub.RawEmail, Password: validPassword}
}

func (s *UserUseCaseShould) TestBuild_WithoutParams_ReturnError() {
	useCase, err := userusecase.NewBuilder().Build()

	assert.Nil(s.T(), useCase)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userusecase.ErrUserRepositoryIsRequired))
	assert.True(s.T(), errors.Is(err, userusecase.ErrTokenUseCaseIsRequired))
	assert.True(s.T(), errors.Is(err, userusecase.ErrPasswordHasherIsRequired))
	assert.True(s.T(), errors.Is(err, userusecase.ErrLogPublisherIsRequired))
}

func (s *UserUseCaseShould) TestRegister_NilData_ReturnError() {
	user, err := s.useCase.Register(s.T().Context(), nil)

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userusecase.ErrDataIsRequired))
}

func (s *UserUseCaseShould) TestRegister_MalformedEmail_ReturnError() {
	data := s.registerData()
	data.Email = "not-an-email"

	user, err := s.useCase.Register(s.T().Context(), data)

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, emailprimitive.ErrEmailInvalidFormat))
}

func (s *UserUseCaseShould) TestRegister_ShortPassword_ReturnError() {
	data := s.registerData()
	data.Password = "short"

	user, err := s.useCase.Register(s.T().Context(), data)

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, userusecase.ErrCodePasswordTooShort))
}

func (s *UserUseCaseShould) TestRegister_EmailAlreadyTaken_ReturnError() {
	s.userRepository.GetByEmailFunc = func(_ context.Context, _ emailprimitive.Email) (*userentity.User, error) {
		return userentitystub.GetValidUser(), nil
	}

	user, err := s.useCase.Register(s.T().Context(), s.registerData())

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, userusecase.ErrCodeEmailAlreadyTaken))
}

func (s *UserUseCaseShould) TestRegister_RepositoryFails_WrapError() {
	cause := errors.New("database is down")
	s.userRepository.InsertFunc = func(context.Context, *userentity.User) error { return cause }

	user, err := s.useCase.Register(s.T().Context(), s.registerData())

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, cause)
	assert.True(s.T(), apperrors.EqualByCode(err, userusecase.ErrCodeOperationFailed))
}

func (s *UserUseCaseShould) TestRegister_ValidData_InsertHashedUser() {
	user, err := s.useCase.Register(s.T().Context(), s.registerData())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), user)
	assert.Equal(s.T(), userentitystub.RawEmail, user.Email().String())
	assert.Equal(s.T(), userusecasemock.StubPasswordHash, user.PasswordHash().String())
	assert.False(s.T(), user.ID().IsZero())
	assert.Same(s.T(), user, s.userRepository.Inserted)
}

func (s *UserUseCaseShould) TestRegister_NeverStoresRawPassword() {
	user, err := s.useCase.Register(s.T().Context(), s.registerData())

	require.NoError(s.T(), err)
	assert.NotEqual(s.T(), validPassword, user.PasswordHash().String())
}

func (s *UserUseCaseShould) TestLogin_UnknownEmail_ReturnInvalidCredentials() {
	loginData := &boundarymodel.LoginData{Email: userentitystub.RawEmail, Password: validPassword}

	result, err := s.useCase.Login(s.T().Context(), loginData)

	assert.Nil(s.T(), result)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userusecase.ErrInvalidCredentials))
}

func (s *UserUseCaseShould) TestLogin_WrongPassword_ReturnInvalidCredentials() {
	s.userRepository.GetByEmailFunc = func(_ context.Context, _ emailprimitive.Email) (*userentity.User, error) {
		return userentitystub.GetValidUser(), nil
	}
	s.passwordHasher.CompareFunc = func(string, string) error { return errors.New("mismatch") }

	loginData := &boundarymodel.LoginData{Email: userentitystub.RawEmail, Password: "wrong password"}

	result, err := s.useCase.Login(s.T().Context(), loginData)

	assert.Nil(s.T(), result)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userusecase.ErrInvalidCredentials))
}

func (s *UserUseCaseShould) TestLogin_ValidCredentials_ReturnToken() {
	s.userRepository.GetByEmailFunc = func(_ context.Context, _ emailprimitive.Email) (*userentity.User, error) {
		return userentitystub.GetValidUser(), nil
	}

	loginData := &boundarymodel.LoginData{Email: userentitystub.RawEmail, Password: validPassword}

	result, err := s.useCase.Login(s.T().Context(), loginData)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	assert.Equal(s.T(), userusecasemock.StubAccessToken, result.AccessToken)
	assert.Equal(s.T(), userusecasemock.StubExpiresIn, result.ExpiresIn)
}

func (s *UserUseCaseShould) TestGetByID_NotFound_ReturnError() {
	user, err := s.useCase.GetByID(s.T().Context(), commonuserentity.NewUserID())

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, userusecase.ErrCodeUserNotFound))
}

func (s *UserUseCaseShould) TestGetByID_ZeroID_ReturnError() {
	user, err := s.useCase.GetByID(s.T().Context(), commonuserentity.UserID{})

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userusecase.ErrDataIsRequired))
}

func (s *UserUseCaseShould) TestGetByID_Found_ReturnUser() {
	expected := userentitystub.GetValidUser()
	s.userRepository.GetByIDFunc = func(context.Context, commonuserentity.UserID) (*userentity.User, error) {
		return expected, nil
	}

	user, err := s.useCase.GetByID(s.T().Context(), expected.ID())

	require.NoError(s.T(), err)
	assert.Same(s.T(), expected, user)
}
