package post_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	postrestcontroller "github.com/BenyaChef/micro/blog/adapter/controller/rest/post"
	authgateway "github.com/BenyaChef/micro/blog/adapter/gateway/auth"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	postusecase "github.com/BenyaChef/micro/blog/domain/usecase/post"
	postusecasemock "github.com/BenyaChef/micro/blog/test/domain/usecase/post/mock"
	commonauthgateway "github.com/BenyaChef/micro/common/auth/adapter/gateway"
	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	"github.com/BenyaChef/micro/infrastructure/grpcclient"
	"github.com/BenyaChef/micro/infrastructure/grpcserver"
	"github.com/BenyaChef/micro/infrastructure/grpcserver/grpcmodel"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
	restservermiddleware "github.com/BenyaChef/micro/infrastructure/restserver/middleware"
	restserverservice "github.com/BenyaChef/micro/infrastructure/restserver/service"
)

const dialDeadline = 3 * time.Second

type fakeAuthService struct {
	authproto.UnimplementedAuthServiceServer

	tokens map[string]string
}

func (f *fakeAuthService) ValidateToken(
	_ context.Context,
	request *authproto.ValidateTokenRequest,
) (*authproto.ValidateTokenResponse, error) {
	userID, known := f.tokens[request.GetAccessToken()]
	if !known {
		return &authproto.ValidateTokenResponse{
			Result: &authproto.ValidateTokenResponse_Error{
				Error: &grpcmodel.Error{Code: "10b754c6-002", Msg: "Token is invalid"},
			},
		}, nil
	}

	return &authproto.ValidateTokenResponse{
		Result: &authproto.ValidateTokenResponse_Token{
			Token: &authproto.ValidatedToken{UserId: userID},
		},
	}, nil
}

func freeTCPPort(t *testing.T) string {
	t.Helper()

	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(t.Context(), "tcp", "localhost:0")
	require.NoError(t, err)

	addr, ok := listener.Addr().(*net.TCPAddr)
	require.True(t, ok)

	port := strconv.Itoa(addr.Port)
	require.NoError(t, listener.Close())

	return port
}

func waitForTCPPort(t *testing.T, port string) {
	t.Helper()

	var dialer net.Dialer

	deadline := time.Now().Add(dialDeadline)
	for time.Now().Before(deadline) {
		conn, err := dialer.DialContext(t.Context(), "tcp", net.JoinHostPort("localhost", port))
		if err == nil {
			require.NoError(t, conn.Close())

			return
		}

		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("поддельный auth не поднялся на порту %s за %s", port, dialDeadline)
}

func startFakeAuth(t *testing.T, tokens map[string]string) string {
	t.Helper()

	logPublisher := loggerteststub.NewLogPublisherStub()
	port := freeTCPPort(t)

	server, err := grpcserver.NewBuilder().LogPublisher(logPublisher).Port(port).Build()
	require.NoError(t, err)

	authproto.RegisterAuthServiceServer(server.Server(), &fakeAuthService{tokens: tokens})

	serverCtx, stopServer := context.WithCancel(context.Background())
	runFailed := make(chan error, 1)

	go func() { runFailed <- server.Run(serverCtx) }()

	t.Cleanup(func() {
		stopServer()
		require.NoError(t, <-runFailed)
	})

	waitForTCPPort(t, port)

	return net.JoinHostPort("localhost", port)
}

func buildBlogRouter(t *testing.T, authAddr string, repository *postusecasemock.PostRepositoryMock) chi.Router {
	t.Helper()

	logPublisher := loggerteststub.NewLogPublisherStub()

	client, err := grpcclient.NewBuilder().Target(authAddr).LogPublisher(logPublisher).Build()
	require.NoError(t, err)

	t.Cleanup(func() { client.Close(context.Background()) })

	commonGateway, err := commonauthgateway.NewBuilder().
		Client(authproto.NewAuthServiceClient(client.Connection())).
		LogPublisher(logPublisher).
		Build()
	require.NoError(t, err)

	blogGateway, err := authgateway.NewBuilder().
		AuthGateway(commonGateway).
		LogPublisher(logPublisher).
		Build()
	require.NoError(t, err)

	tokenValidator, err := authgateway.NewTokenValidator(blogGateway)
	require.NoError(t, err)

	useCase, err := postusecase.NewBuilder().
		PostRepository(repository).
		LogPublisher(logPublisher).
		Build()
	require.NoError(t, err)

	errorResponseService, err := restserverservice.NewErrorResponseService(
		postrestcontroller.NewErrorResolver(), logPublisher,
	)
	require.NoError(t, err)

	responseService, err := restserverservice.NewResponseService(errorResponseService, logPublisher)
	require.NoError(t, err)

	baseController, err := restservercontroller.NewBaseController(responseService, logPublisher)
	require.NoError(t, err)

	controller, err := postrestcontroller.NewBuilder().
		BaseController(baseController).
		PostUseCase(useCase).
		Build()
	require.NoError(t, err)

	router := chi.NewRouter()
	router.Route("/posts", func(postsRouter chi.Router) {
		postsRouter.Get("/", controller.List)
		postsRouter.Get("/{postID}", controller.GetByID)

		postsRouter.Group(func(protected chi.Router) {
			protected.Use(restservermiddleware.Auth(tokenValidator, responseService))

			protected.Post("/", controller.Create)
			protected.Put("/{postID}", controller.Update)
			protected.Delete("/{postID}", controller.Delete)
		})
	})

	return router
}

func TestBlogToAuthOverRealSocket(t *testing.T) {
	t.Parallel()

	alice := commonuserentity.NewUserID()
	bob := commonuserentity.NewUserID()

	authAddr := startFakeAuth(t, map[string]string{
		"alice-token": alice.String(),
		"bob-token":   bob.String(),
	})

	repository := postusecasemock.NewPostRepositoryMock()
	router := buildBlogRouter(t, authAddr, repository)

	call := func(method, target, body, token string) *httptest.ResponseRecorder {
		t.Helper()

		request := httptest.NewRequestWithContext(t.Context(), method, target, strings.NewReader(body))
		if token != "" {
			request.Header.Set("Authorization", "Bearer "+token)
		}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		return recorder
	}

	var alicePostID string

	t.Run("без токена 401", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, call(http.MethodPost, "/posts", validBody, "").Code)
	})

	t.Run("битый токен отвергается настоящим auth 401", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, call(http.MethodPost, "/posts", validBody, "garbage").Code)
	})

	t.Run("алиса создаёт статью", func(t *testing.T) {
		recorder := call(http.MethodPost, "/posts", validBody, "alice-token")
		require.Equal(t, http.StatusCreated, recorder.Code)

		var body map[string]any
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
		assert.Equal(t, alice.String(), body["author_id"])

		alicePostID, _ = body["post_id"].(string)
		require.NotEmpty(t, alicePostID)
	})

	t.Run("боб не может править чужую статью 403", func(t *testing.T) {
		repository.GetByIDFunc = func(context.Context, postentity.PostID) (*postentity.Post, error) {
			return repository.Inserted, nil
		}

		recorder := call(http.MethodPut, "/posts/"+alicePostID, validBody, "bob-token")

		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Nil(t, repository.Updated)
	})

	t.Run("алиса правит свою статью", func(t *testing.T) {
		recorder := call(http.MethodPut, "/posts/"+alicePostID, validBody, "alice-token")

		assert.Equal(t, http.StatusOK, recorder.Code)
		require.NotNil(t, repository.Updated)
	})

	t.Run("боб не может удалить чужую статью 403", func(t *testing.T) {
		assert.Equal(t, http.StatusForbidden, call(http.MethodDelete, "/posts/"+alicePostID, "", "bob-token").Code)
	})

	t.Run("чтение списка публично", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, call(http.MethodGet, "/posts", "", "").Code)
	})
}
