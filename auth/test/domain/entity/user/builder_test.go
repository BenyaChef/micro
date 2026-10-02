package user_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	userentitystub "github.com/BenyaChef/micro/auth/test/domain/entity/user/stub"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type BuilderShould struct {
	suite.Suite
}

func TestBuilderShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuilderShould))
}

func (s *BuilderShould) TestBuild_WithoutParams_ReturnError() {
	user, err := userentity.NewBuilder().Build()

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userentity.ErrUserIDIsRequired))
	assert.True(s.T(), errors.Is(err, userentity.ErrEmailIsRequired))
	assert.True(s.T(), errors.Is(err, userentity.ErrPasswordHashIsRequired))
}

func (s *BuilderShould) TestBuild_WithoutUserID_ReturnOnlyUserIDError() {
	user, err := userentity.NewBuilder().
		Email(userentitystub.GetEmail()).
		PasswordHash(userentitystub.GetPasswordHash()).
		Build()

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userentity.ErrUserIDIsRequired))
	assert.False(s.T(), errors.Is(err, userentity.ErrEmailIsRequired))
	assert.False(s.T(), errors.Is(err, userentity.ErrPasswordHashIsRequired))
}

func (s *BuilderShould) TestBuild_WithoutEmail_ReturnOnlyEmailError() {
	user, err := userentity.NewBuilder().
		ID(userentitystub.GetUserID()).
		PasswordHash(userentitystub.GetPasswordHash()).
		Build()

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userentity.ErrEmailIsRequired))
	assert.False(s.T(), errors.Is(err, userentity.ErrUserIDIsRequired))
	assert.False(s.T(), errors.Is(err, userentity.ErrPasswordHashIsRequired))
}

func (s *BuilderShould) TestBuild_WithoutPasswordHash_ReturnOnlyPasswordHashError() {
	user, err := userentity.NewBuilder().
		ID(userentitystub.GetUserID()).
		Email(userentitystub.GetEmail()).
		Build()

	assert.Nil(s.T(), user)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userentity.ErrPasswordHashIsRequired))
	assert.False(s.T(), errors.Is(err, userentity.ErrUserIDIsRequired))
	assert.False(s.T(), errors.Is(err, userentity.ErrEmailIsRequired))
}

func (s *BuilderShould) TestBuild_AllParams_ReturnUser() {
	userID := userentitystub.GetUserID()
	email := userentitystub.GetEmail()
	passwordHash := userentitystub.GetPasswordHash()

	user, err := userentity.NewBuilder().
		ID(userID).
		Email(email).
		PasswordHash(passwordHash).
		Build()

	require.NoError(s.T(), err)
	require.NotNil(s.T(), user)
	assert.Equal(s.T(), userID, user.ID())
	assert.Equal(s.T(), email, user.Email())
	assert.Equal(s.T(), passwordHash, user.PasswordHash())
}

func (s *BuilderShould) TestBuild_WithoutTimestamps_FillThemWithNow() {
	user, err := userentitystub.GetUserBuilder().Build()

	require.NoError(s.T(), err)
	require.NotNil(s.T(), user)
	assert.False(s.T(), user.CreatedAt().IsZero())
	assert.False(s.T(), user.UpdatedAt().IsZero())
	assert.Equal(s.T(), user.CreatedAt(), user.UpdatedAt())
}

func (s *BuilderShould) TestBuild_WithTimestamps_KeepThem() {
	createdAt := timeprimitive.Now()
	updatedAt := timeprimitive.Now()

	user, err := userentitystub.GetUserBuilder().
		CreatedAt(createdAt).
		UpdatedAt(updatedAt).
		Build()

	require.NoError(s.T(), err)
	require.NotNil(s.T(), user)
	assert.Equal(s.T(), createdAt, user.CreatedAt())
	assert.Equal(s.T(), updatedAt, user.UpdatedAt())
}
