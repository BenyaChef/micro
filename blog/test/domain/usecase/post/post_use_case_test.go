package post_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	boundarymodel "github.com/BenyaChef/micro/blog/boundary/model"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	postusecase "github.com/BenyaChef/micro/blog/domain/usecase/post"
	postentitystub "github.com/BenyaChef/micro/blog/test/domain/entity/post/stub"
	postusecasemock "github.com/BenyaChef/micro/blog/test/domain/usecase/post/mock"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
)

type PostUseCaseShould struct {
	suite.Suite

	postRepository *postusecasemock.PostRepositoryMock
	useCase        *postusecase.PostUseCase
}

func TestPostUseCaseShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PostUseCaseShould))
}

func (s *PostUseCaseShould) SetupTest() {
	s.postRepository = postusecasemock.NewPostRepositoryMock()

	useCase, err := postusecase.NewBuilder().
		PostRepository(s.postRepository).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.useCase = useCase
}

func (s *PostUseCaseShould) TestBuild_WithoutParams_ReturnError() {
	useCase, err := postusecase.NewBuilder().Build()

	assert.Nil(s.T(), useCase)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postusecase.ErrPostRepositoryIsRequired))
	assert.True(s.T(), errors.Is(err, postusecase.ErrLogPublisherIsRequired))
}

func (s *PostUseCaseShould) TestCreate_NilInput_ReturnError() {
	post, err := s.useCase.Create(s.T().Context(), nil)

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postusecase.ErrInputIsRequired))
}

func (s *PostUseCaseShould) TestCreate_ZeroAuthor_ReturnError() {
	input := &boundarymodel.CreatePostInput{Title: postentitystub.RawTitle, Body: postentitystub.RawBody}

	post, err := s.useCase.Create(s.T().Context(), input)

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postusecase.ErrAuthorIDIsRequired))
}

func (s *PostUseCaseShould) TestCreate_EmptyTitle_ReturnError() {
	input := &boundarymodel.CreatePostInput{
		AuthorID: commonuserentity.NewUserID(),
		Title:    "   ",
		Body:     postentitystub.RawBody,
	}

	post, err := s.useCase.Create(s.T().Context(), input)

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postentity.ErrTitleIsEmpty))
}

func (s *PostUseCaseShould) TestCreate_ValidInput_InsertPost() {
	authorID := commonuserentity.NewUserID()
	input := &boundarymodel.CreatePostInput{
		AuthorID: authorID,
		Title:    postentitystub.RawTitle,
		Body:     postentitystub.RawBody,
	}

	post, err := s.useCase.Create(s.T().Context(), input)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), post)
	assert.Equal(s.T(), authorID, post.AuthorID())
	assert.Same(s.T(), post, s.postRepository.Inserted)
}

func (s *PostUseCaseShould) TestGetByID_NotFound_ReturnNotFoundError() {
	post, err := s.useCase.GetByID(s.T().Context(), postentity.NewPostID())

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postusecase.ErrCodePostNotFound))
}

func (s *PostUseCaseShould) TestList_ZeroLimit_UsesDefault() {
	_, err := s.useCase.List(s.T().Context(), &boundarymodel.ListPostsInput{})

	require.NoError(s.T(), err)
	assert.Equal(s.T(), postusecase.DefaultListLimit, s.postRepository.CalledLimit)
	assert.Equal(s.T(), 0, s.postRepository.CalledOffset)
}

func (s *PostUseCaseShould) TestList_HugeLimit_IsCapped() {
	_, err := s.useCase.List(s.T().Context(), &boundarymodel.ListPostsInput{Limit: 10_000, Offset: -5})

	require.NoError(s.T(), err)
	assert.Equal(s.T(), postusecase.MaxListLimit, s.postRepository.CalledLimit)
	assert.Equal(s.T(), 0, s.postRepository.CalledOffset)
}

func (s *PostUseCaseShould) TestUpdate_ForeignPost_ReturnNotAuthorError() {
	owner := commonuserentity.NewUserID()
	existing := postentitystub.GetPostAuthoredBy(owner)
	s.postRepository.GetByIDFunc = func(context.Context, postentity.PostID) (*postentity.Post, error) {
		return existing, nil
	}

	input := &boundarymodel.UpdatePostInput{
		PostID:   existing.ID(),
		AuthorID: commonuserentity.NewUserID(),
		Title:    "чужая правка",
		Body:     "чужое тело",
	}

	post, err := s.useCase.Update(s.T().Context(), input)

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postusecase.ErrCodeNotAuthor))
	assert.Nil(s.T(), s.postRepository.Updated)
}

func (s *PostUseCaseShould) TestUpdate_OwnPost_KeepsIDAndCreatedAt() {
	owner := commonuserentity.NewUserID()
	existing := postentitystub.GetPostAuthoredBy(owner)
	s.postRepository.GetByIDFunc = func(context.Context, postentity.PostID) (*postentity.Post, error) {
		return existing, nil
	}

	input := &boundarymodel.UpdatePostInput{
		PostID:   existing.ID(),
		AuthorID: owner,
		Title:    "новый заголовок",
		Body:     "новое тело",
	}

	post, err := s.useCase.Update(s.T().Context(), input)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), post)
	assert.Equal(s.T(), existing.ID(), post.ID())
	assert.Equal(s.T(), existing.CreatedAt(), post.CreatedAt())
	assert.Equal(s.T(), "новый заголовок", post.Title().String())
	assert.Same(s.T(), post, s.postRepository.Updated)
}

func (s *PostUseCaseShould) TestDelete_ForeignPost_ReturnNotAuthorError() {
	owner := commonuserentity.NewUserID()
	existing := postentitystub.GetPostAuthoredBy(owner)
	s.postRepository.GetByIDFunc = func(context.Context, postentity.PostID) (*postentity.Post, error) {
		return existing, nil
	}

	err := s.useCase.Delete(s.T().Context(), existing.ID(), commonuserentity.NewUserID())

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postusecase.ErrCodeNotAuthor))
	assert.True(s.T(), s.postRepository.Deleted.IsZero())
}

func (s *PostUseCaseShould) TestDelete_OwnPost_DeletesIt() {
	owner := commonuserentity.NewUserID()
	existing := postentitystub.GetPostAuthoredBy(owner)
	s.postRepository.GetByIDFunc = func(context.Context, postentity.PostID) (*postentity.Post, error) {
		return existing, nil
	}

	err := s.useCase.Delete(s.T().Context(), existing.ID(), owner)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), existing.ID(), s.postRepository.Deleted)
}

func (s *PostUseCaseShould) TestDelete_MissingPost_ReturnNotFound() {
	err := s.useCase.Delete(s.T().Context(), postentity.NewPostID(), commonuserentity.NewUserID())

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postusecase.ErrCodePostNotFound))
}

func (s *PostUseCaseShould) TestCreate_RepositoryFails_WrapError() {
	cause := errors.New("database is down")
	s.postRepository.InsertFunc = func(context.Context, *postentity.Post) error { return cause }

	input := &boundarymodel.CreatePostInput{
		AuthorID: commonuserentity.NewUserID(),
		Title:    postentitystub.RawTitle,
		Body:     postentitystub.RawBody,
	}

	post, err := s.useCase.Create(s.T().Context(), input)

	assert.Nil(s.T(), post)
	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, cause)
	assert.True(s.T(), apperrors.EqualByCode(err, postusecase.ErrCodeOperationFailed))
}
