package post_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	postentitystub "github.com/BenyaChef/micro/blog/test/domain/entity/post/stub"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
)

type BuilderShould struct {
	suite.Suite
}

func TestBuilderShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuilderShould))
}

func (s *BuilderShould) TestBuild_WithoutParams_ReturnAllErrors() {
	post, err := postentity.NewBuilder().Build()

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postentity.ErrPostIDIsRequired))
	assert.True(s.T(), errors.Is(err, postentity.ErrAuthorIDIsRequired))
	assert.True(s.T(), errors.Is(err, postentity.ErrTitleIsRequired))
	assert.True(s.T(), errors.Is(err, postentity.ErrBodyIsRequired))
}

func (s *BuilderShould) TestBuild_WithoutAuthorID_ReturnOnlyAuthorError() {
	post, err := postentity.NewBuilder().
		ID(postentity.NewPostID()).
		Title(postentitystub.GetTitle()).
		Body(postentitystub.GetBody()).
		Build()

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postentity.ErrAuthorIDIsRequired))
	assert.False(s.T(), errors.Is(err, postentity.ErrPostIDIsRequired))
}

func (s *BuilderShould) TestBuild_AllParams_ReturnPost() {
	authorID := commonuserentity.NewUserID()

	post := postentitystub.GetPostAuthoredBy(authorID)

	assert.Equal(s.T(), authorID, post.AuthorID())
	assert.Equal(s.T(), postentitystub.RawTitle, post.Title().String())
	assert.False(s.T(), post.ID().IsZero())
}

func (s *BuilderShould) TestBuild_WithoutTimestamps_FillThemWithNow() {
	post := postentitystub.GetValidPost()

	assert.False(s.T(), post.CreatedAt().IsZero())
	assert.False(s.T(), post.UpdatedAt().IsZero())
	assert.Equal(s.T(), post.CreatedAt(), post.UpdatedAt())
}

func (s *BuilderShould) TestIsAuthoredBy_OwnAuthor_ReturnTrue() {
	authorID := commonuserentity.NewUserID()
	post := postentitystub.GetPostAuthoredBy(authorID)

	assert.True(s.T(), post.IsAuthoredBy(authorID))
	assert.False(s.T(), post.IsAuthoredBy(commonuserentity.NewUserID()))
}

func (s *BuilderShould) TestTitleFrom_EmptyOrBlank_ReturnError() {
	for _, raw := range []string{"", "   ", "\n\t"} {
		title, err := postentity.TitleFrom(raw)

		require.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, postentity.ErrTitleIsEmpty))
		assert.True(s.T(), title.IsZero())
	}
}

func (s *BuilderShould) TestTitleFrom_TooLong_ReturnError() {
	title, err := postentity.TitleFrom(strings.Repeat("я", postentity.MaxTitleLength+1))

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postentity.ErrCodeTitleTooLong))
	assert.True(s.T(), title.IsZero())
}

func (s *BuilderShould) TestTitleFrom_Padded_IsTrimmed() {
	title, err := postentity.TitleFrom("  заголовок  ")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), "заголовок", title.String())
}

func (s *BuilderShould) TestBodyFrom_EmptyOrTooLong_ReturnError() {
	body, err := postentity.BodyFrom("  ")
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postentity.ErrBodyIsEmpty))
	assert.True(s.T(), body.IsZero())

	body, err = postentity.BodyFrom(strings.Repeat("a", postentity.MaxBodyLength+1))
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postentity.ErrCodeBodyTooLong))
	assert.True(s.T(), body.IsZero())
}

func (s *BuilderShould) TestPostIDFrom_InvalidUUID_ReturnError() {
	for _, raw := range []string{"", "not-a-uuid", "550e8400-e29b"} {
		postID, err := postentity.PostIDFrom(raw)

		require.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, postentity.ErrPostIDInvalidFormat))
		assert.True(s.T(), postID.IsZero())
	}
}
