package post_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	postrepository "github.com/BenyaChef/micro/blog/adapter/repository/post"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	postentitystub "github.com/BenyaChef/micro/blog/test/domain/entity/post/stub"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	"github.com/BenyaChef/micro/infrastructure/postgres"
)

const envTestDSN = "POSTGRES_TEST_DSN"

type PostRepositoryShould struct {
	suite.Suite

	client     *postgres.Client
	repository *postrepository.PostRepository
}

func TestPostRepositoryShould(t *testing.T) {
	suite.Run(t, new(PostRepositoryShould))
}

func (s *PostRepositoryShould) SetupSuite() {
	dsn := os.Getenv(envTestDSN)
	if dsn == "" {
		s.T().Skipf("%s is not set, skipping test against live database", envTestDSN)
	}

	client, err := postgres.NewBuilder().DSN(dsn).Build(s.T().Context())
	require.NoError(s.T(), err)

	s.client = client

	repository, err := postrepository.NewBuilder().
		Executor(client).
		LogPublisher(loggerteststub.NewLogPublisherStub()).
		Build()
	require.NoError(s.T(), err)

	s.repository = repository
}

func (s *PostRepositoryShould) TearDownSuite() {
	if s.client != nil {
		s.client.Close()
	}
}

func (s *PostRepositoryShould) insertFreshPost(authorID commonuserentity.UserID) *postentity.Post {
	s.T().Helper()

	post := postentitystub.GetPostAuthoredBy(authorID)
	require.NoError(s.T(), s.repository.Insert(s.T().Context(), post))

	s.T().Cleanup(func() {
		_, err := s.client.Exec(context.Background(), "DELETE FROM posts WHERE id = $1", post.ID().String())
		require.NoError(s.T(), err)
	})

	return post
}

func (s *PostRepositoryShould) TestBuild_WithoutParams_ReturnError() {
	repository, err := postrepository.NewBuilder().Build()

	assert.Nil(s.T(), repository)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postrepository.ErrExecutorIsRequired))
	assert.True(s.T(), errors.Is(err, postrepository.ErrLogPublisherIsRequired))
}

func (s *PostRepositoryShould) TestInsert_NilPost_ReturnError() {
	err := s.repository.Insert(s.T().Context(), nil)

	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postrepository.ErrPostIsRequired))
}

func (s *PostRepositoryShould) TestGetByID_InsertedPost_ReturnEqualEntity() {
	inserted := s.insertFreshPost(commonuserentity.NewUserID())

	found, err := s.repository.GetByID(s.T().Context(), inserted.ID())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), inserted.ID(), found.ID())
	assert.Equal(s.T(), inserted.AuthorID(), found.AuthorID())
	assert.Equal(s.T(), inserted.Title(), found.Title())
	assert.Equal(s.T(), inserted.Body(), found.Body())
	assert.Equal(s.T(), inserted.CreatedAt(), found.CreatedAt())
}

func (s *PostRepositoryShould) TestGetByID_UnknownPost_ReturnNilNil() {
	found, err := s.repository.GetByID(s.T().Context(), postentity.NewPostID())

	require.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *PostRepositoryShould) TestList_ReturnsNewestFirst() {
	authorID := commonuserentity.NewUserID()
	older := s.insertFreshPost(authorID)
	newer := s.insertFreshPost(authorID)

	posts, err := s.repository.List(s.T().Context(), 100, 0)

	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), posts)

	positions := map[string]int{}
	for i, post := range posts {
		positions[post.ID().String()] = i
	}

	require.Contains(s.T(), positions, older.ID().String())
	require.Contains(s.T(), positions, newer.ID().String())
	assert.Less(s.T(), positions[newer.ID().String()], positions[older.ID().String()])
}

func (s *PostRepositoryShould) TestList_LimitIsRespected() {
	authorID := commonuserentity.NewUserID()
	s.insertFreshPost(authorID)
	s.insertFreshPost(authorID)

	posts, err := s.repository.List(s.T().Context(), 1, 0)

	require.NoError(s.T(), err)
	assert.Len(s.T(), posts, 1)
}

func (s *PostRepositoryShould) TestUpdate_ChangesTitleAndBody() {
	inserted := s.insertFreshPost(commonuserentity.NewUserID())

	newTitle, err := postentity.TitleFrom("обновлённый заголовок")
	require.NoError(s.T(), err)

	newBody, err := postentity.BodyFrom("обновлённое тело")
	require.NoError(s.T(), err)

	updated, err := postentity.NewBuilder().
		ID(inserted.ID()).
		AuthorID(inserted.AuthorID()).
		Title(newTitle).
		Body(newBody).
		CreatedAt(inserted.CreatedAt()).
		Build()
	require.NoError(s.T(), err)

	require.NoError(s.T(), s.repository.Update(s.T().Context(), updated))

	found, err := s.repository.GetByID(s.T().Context(), inserted.ID())
	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), "обновлённый заголовок", found.Title().String())
	assert.Equal(s.T(), inserted.CreatedAt(), found.CreatedAt())
}

func (s *PostRepositoryShould) TestUpdate_MissingPost_ReturnVanishedError() {
	ghost := postentitystub.GetValidPost()

	err := s.repository.Update(s.T().Context(), ghost)

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postrepository.ErrCodePostVanished))
}

func (s *PostRepositoryShould) TestDelete_RemovesPost() {
	inserted := s.insertFreshPost(commonuserentity.NewUserID())

	require.NoError(s.T(), s.repository.Delete(s.T().Context(), inserted.ID()))

	found, err := s.repository.GetByID(s.T().Context(), inserted.ID())
	require.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *PostRepositoryShould) TestDelete_MissingPost_ReturnVanishedError() {
	err := s.repository.Delete(s.T().Context(), postentity.NewPostID())

	require.Error(s.T(), err)
	assert.True(s.T(), apperrors.EqualByCode(err, postrepository.ErrCodePostVanished))
}
