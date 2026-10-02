package user_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	userentitystub "github.com/BenyaChef/micro/auth/test/domain/entity/user/stub"
)

type PasswordHashShould struct {
	suite.Suite
}

func TestPasswordHashShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PasswordHashShould))
}

func (s *PasswordHashShould) TestPasswordHashFrom_EmptyValue_ReturnError() {
	passwordHash, err := userentity.PasswordHashFrom("")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, userentity.ErrPasswordHashIsEmpty))
	assert.True(s.T(), passwordHash.IsZero())
}

func (s *PasswordHashShould) TestPasswordHashFrom_ValidValue_ReturnHash() {
	passwordHash, err := userentity.PasswordHashFrom(userentitystub.RawPasswordHash)

	require.NoError(s.T(), err)
	assert.False(s.T(), passwordHash.IsZero())
	assert.Equal(s.T(), userentitystub.RawPasswordHash, passwordHash.String())
}
