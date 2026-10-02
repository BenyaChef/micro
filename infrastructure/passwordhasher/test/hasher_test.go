package passwordhasher_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	"github.com/BenyaChef/micro/infrastructure/passwordhasher"
)

const testPassword = "correct horse battery staple"

type HasherShould struct {
	suite.Suite
}

func TestHasherShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HasherShould))
}

func (s *HasherShould) newHasher() *passwordhasher.Hasher {
	hasher, err := passwordhasher.NewBuilder().Cost(bcrypt.MinCost).Build()
	require.NoError(s.T(), err)

	return hasher
}

func (s *HasherShould) TestBuild_CostOutOfRange_ReturnError() {
	hasher, err := passwordhasher.NewBuilder().Cost(bcrypt.MaxCost + 1).Build()

	assert.Nil(s.T(), hasher)
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, passwordhasher.ErrCodeCostOutOfRange))
}

func (s *HasherShould) TestHash_EmptyPassword_ReturnError() {
	passwordHash, err := s.newHasher().Hash("")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, passwordhasher.ErrPasswordIsEmpty))
	assert.Empty(s.T(), passwordHash)
}

func (s *HasherShould) TestHash_SamePasswordTwice_ReturnDifferentHashes() {
	hasher := s.newHasher()

	first, err := hasher.Hash(testPassword)
	require.NoError(s.T(), err)

	second, err := hasher.Hash(testPassword)
	require.NoError(s.T(), err)

	assert.NotEqual(s.T(), first, second)
}

func (s *HasherShould) TestCompare_CorrectPassword_ReturnNil() {
	hasher := s.newHasher()

	passwordHash, err := hasher.Hash(testPassword)
	require.NoError(s.T(), err)

	require.NoError(s.T(), hasher.Compare(passwordHash, testPassword))
}

func (s *HasherShould) TestCompare_WrongPassword_ReturnMismatch() {
	hasher := s.newHasher()

	passwordHash, err := hasher.Hash(testPassword)
	require.NoError(s.T(), err)

	err = hasher.Compare(passwordHash, "wrong password")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, passwordhasher.ErrPasswordMismatch))
}

func (s *HasherShould) TestCompare_EmptyPassword_ReturnError() {
	err := s.newHasher().Compare("$2a$04$irrelevant", "")

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, passwordhasher.ErrPasswordIsEmpty))
}
