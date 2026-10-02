package user_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

type UserIDShould struct {
	suite.Suite
}

func TestUserIDShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserIDShould))
}

func (s *UserIDShould) TestNewUserID_GeneratesUniqueID() {
	first := commonuserentity.NewUserID()
	second := commonuserentity.NewUserID()

	assert.False(s.T(), first.IsZero())
	assert.False(s.T(), second.IsZero())
	assert.NotEqual(s.T(), first.String(), second.String())
}

func (s *UserIDShould) TestUserIDFrom_ValidUUID_ReturnUserID() {
	rawID := "550e8400-e29b-41d4-a716-446655440000"

	userID, err := commonuserentity.UserIDFrom(rawID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), rawID, userID.String())
}

func (s *UserIDShould) TestUserIDFrom_InvalidUUID_ReturnError() {
	invalidCases := []struct {
		name  string
		rawID string
	}{
		{"empty", ""},
		{"not uuid", "not-a-uuid"},
		{"too short", "550e8400-e29b-41d4"},
	}

	for _, testCase := range invalidCases {
		s.Run(testCase.name, func() {
			userID, err := commonuserentity.UserIDFrom(testCase.rawID)

			require.Error(s.T(), err)
			assert.True(s.T(), errors.Is(err, commonuserentity.ErrUserIDInvalidFormat))
			assert.True(s.T(), userID.IsZero())
		})
	}
}
