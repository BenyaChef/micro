package post_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	postrestcontroller "github.com/BenyaChef/micro/blog/adapter/controller/rest/post"
	authgateway "github.com/BenyaChef/micro/blog/adapter/gateway/auth"
	boundarymodel "github.com/BenyaChef/micro/blog/boundary/model"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	postusecase "github.com/BenyaChef/micro/blog/domain/usecase/post"
	postrestcontrollermock "github.com/BenyaChef/micro/blog/test/adapter/controller/rest/post/mock"
	postentitystub "github.com/BenyaChef/micro/blog/test/domain/entity/post/stub"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
	restservermiddleware "github.com/BenyaChef/micro/infrastructure/restserver/middleware"
	restserverservice "github.com/BenyaChef/micro/infrastructure/restserver/service"
)

const validBody = `{"title":"заголовок","body":"тело статьи"}`

type PostControllerShould struct {
	suite.Suite

	postUseCase    *postrestcontrollermock.PostUseCaseMock
	tokenValidator *postrestcontrollermock.TokenValidatorMock
	callerID       commonuserentity.UserID
	router         chi.Router
}

func TestPostControllerShould(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PostControllerShould))
}

func (s *PostControllerShould) SetupTest() {
	logPublisher := loggerteststub.NewLogPublisherStub()

	errorResponseService, err := restserverservice.NewErrorResponseService(
		postrestcontroller.NewErrorResolver(), logPublisher,
	)
	require.NoError(s.T(), err)

	responseService, err := restserverservice.NewResponseService(errorResponseService, logPublisher)
	require.NoError(s.T(), err)

	baseController, err := restservercontroller.NewBaseController(responseService, logPublisher)
	require.NoError(s.T(), err)

	s.postUseCase = postrestcontrollermock.NewPostUseCaseMock()
	s.callerID = commonuserentity.NewUserID()
	s.tokenValidator = postrestcontrollermock.NewTokenValidatorMock(s.callerID.String())

	controller, err := postrestcontroller.NewBuilder().
		BaseController(baseController).
		PostUseCase(s.postUseCase).
		Build()
	require.NoError(s.T(), err)

	router := chi.NewRouter()
	router.Route("/posts", func(postsRouter chi.Router) {
		postsRouter.Get("/", controller.List)
		postsRouter.Get("/{postID}", controller.GetByID)

		postsRouter.Group(func(protected chi.Router) {
			protected.Use(restservermiddleware.Auth(s.tokenValidator, responseService))

			protected.Post("/", controller.Create)
			protected.Put("/{postID}", controller.Update)
			protected.Delete("/{postID}", controller.Delete)
		})
	})

	s.router = router
}

func (s *PostControllerShould) call(method, target, body, token string) *httptest.ResponseRecorder {
	s.T().Helper()

	request := httptest.NewRequestWithContext(s.T().Context(), method, target, strings.NewReader(body))
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, request)

	return recorder
}

func (s *PostControllerShould) TestBuild_WithoutParams_ReturnError() {
	controller, err := postrestcontroller.NewBuilder().Build()

	assert.Nil(s.T(), controller)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, postrestcontroller.ErrBaseControllerIsRequired))
	assert.True(s.T(), errors.Is(err, postrestcontroller.ErrPostUseCaseIsRequired))
}

func (s *PostControllerShould) TestCreate_WithoutToken_Return401() {
	recorder := s.call(http.MethodPost, "/posts", validBody, "")

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *PostControllerShould) TestUpdate_WithoutToken_Return401() {
	target := "/posts/" + postentity.NewPostID().String()

	recorder := s.call(http.MethodPut, target, validBody, "")

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *PostControllerShould) TestDelete_WithoutToken_Return401() {
	target := "/posts/" + postentity.NewPostID().String()

	recorder := s.call(http.MethodDelete, target, "", "")

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *PostControllerShould) TestCreate_RejectedToken_Return401() {
	s.tokenValidator.ValidateTokenFunc = func(context.Context, string) (string, error) {
		return "", authgateway.ErrTokenRejected(errors.New("token is expired"))
	}

	recorder := s.call(http.MethodPost, "/posts", validBody, "stale")

	assert.Equal(s.T(), http.StatusUnauthorized, recorder.Code)
}

func (s *PostControllerShould) TestCreate_AuthServiceDown_Return503() {
	s.tokenValidator.ValidateTokenFunc = func(context.Context, string) (string, error) {
		return "", authgateway.ErrAuthUnavailable(errors.New("connection refused"))
	}

	recorder := s.call(http.MethodPost, "/posts", validBody, "good")

	assert.Equal(s.T(), http.StatusServiceUnavailable, recorder.Code)
}

func (s *PostControllerShould) TestUpdate_ForeignPost_Return403() {
	postID := postentity.NewPostID()
	s.postUseCase.UpdateFunc = func(context.Context, *boundarymodel.UpdatePostInput) (*postentity.Post, error) {
		return nil, postusecase.ErrNotAuthor(postID.String(), s.callerID.String())
	}

	recorder := s.call(http.MethodPut, "/posts/"+postID.String(), validBody, "good")

	assert.Equal(s.T(), http.StatusForbidden, recorder.Code)
}

func (s *PostControllerShould) TestDelete_ForeignPost_Return403() {
	postID := postentity.NewPostID()
	s.postUseCase.DeleteFunc = func(context.Context, postentity.PostID, commonuserentity.UserID) error {
		return postusecase.ErrNotAuthor(postID.String(), s.callerID.String())
	}

	recorder := s.call(http.MethodDelete, "/posts/"+postID.String(), "", "good")

	assert.Equal(s.T(), http.StatusForbidden, recorder.Code)
}

func (s *PostControllerShould) TestGetByID_MissingPost_Return404() {
	postID := postentity.NewPostID()
	s.postUseCase.GetByIDFunc = func(context.Context, postentity.PostID) (*postentity.Post, error) {
		return nil, postusecase.ErrPostNotFound(postID.String())
	}

	recorder := s.call(http.MethodGet, "/posts/"+postID.String(), "", "")

	assert.Equal(s.T(), http.StatusNotFound, recorder.Code)
}

func (s *PostControllerShould) TestGetByID_MalformedID_Return400() {
	recorder := s.call(http.MethodGet, "/posts/not-a-uuid", "", "")

	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

func (s *PostControllerShould) TestList_IsPublic_Return200() {
	s.postUseCase.ListFunc = func(context.Context, *boundarymodel.ListPostsInput) ([]*postentity.Post, error) {
		return []*postentity.Post{postentitystub.GetValidPost()}, nil
	}

	recorder := s.call(http.MethodGet, "/posts", "", "")

	require.Equal(s.T(), http.StatusOK, recorder.Code)

	var body []map[string]any
	require.NoError(s.T(), json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Len(s.T(), body, 1)
}

func (s *PostControllerShould) TestList_BadLimit_Return400() {
	recorder := s.call(http.MethodGet, "/posts?limit=abc", "", "")

	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

func (s *PostControllerShould) TestList_PagingFromQuery_ReachesUseCase() {
	recorder := s.call(http.MethodGet, "/posts?limit=5&offset=10", "", "")

	require.Equal(s.T(), http.StatusOK, recorder.Code)
	require.NotNil(s.T(), s.postUseCase.ListInput)
	assert.Equal(s.T(), 5, s.postUseCase.ListInput.Limit)
	assert.Equal(s.T(), 10, s.postUseCase.ListInput.Offset)
}

func (s *PostControllerShould) TestCreate_ValidToken_PassesCallerAsAuthor() {
	expected := postentitystub.GetPostAuthoredBy(s.callerID)
	s.postUseCase.CreateFunc = func(context.Context, *boundarymodel.CreatePostInput) (*postentity.Post, error) {
		return expected, nil
	}

	recorder := s.call(http.MethodPost, "/posts", validBody, "good")

	require.Equal(s.T(), http.StatusCreated, recorder.Code)
	require.NotNil(s.T(), s.postUseCase.CreateInput)
	assert.Equal(s.T(), s.callerID, s.postUseCase.CreateInput.AuthorID)

	var body map[string]any
	require.NoError(s.T(), json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(s.T(), expected.ID().String(), body["post_id"])
}

func (s *PostControllerShould) TestCreate_MissingTitle_Return400() {
	recorder := s.call(http.MethodPost, "/posts", `{"body":"только тело"}`, "good")

	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

func (s *PostControllerShould) TestDelete_OwnPost_Return204() {
	recorder := s.call(http.MethodDelete, "/posts/"+postentity.NewPostID().String(), "", "good")

	assert.Equal(s.T(), http.StatusNoContent, recorder.Code)
	assert.Empty(s.T(), recorder.Body.String())
}
