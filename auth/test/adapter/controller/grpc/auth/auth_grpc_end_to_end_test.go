package auth_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authgrpccontroller "github.com/BenyaChef/micro/auth/adapter/controller/grpc/auth"
	tokenusecase "github.com/BenyaChef/micro/auth/domain/usecase/token"
	commonauthgateway "github.com/BenyaChef/micro/common/auth/adapter/gateway"
	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	"github.com/BenyaChef/micro/infrastructure/grpcserver"
	"github.com/BenyaChef/micro/infrastructure/jwt"
	loggerteststub "github.com/BenyaChef/micro/infrastructure/logger/test/stub"
)

const (
	endToEndSecret = "0123456789abcdef0123456789abcdef"
	dialDeadline   = 3 * time.Second
)

func freePort(t *testing.T) string {
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

func waitForPort(t *testing.T, port string) {
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

	t.Fatalf("gRPC-сервер не поднялся на порту %s за %s", port, dialDeadline)
}

func TestAuthGRPCEndToEnd(t *testing.T) {
	t.Parallel()

	logPublisher := loggerteststub.NewLogPublisherStub()
	port := freePort(t)

	tokenService, err := jwt.NewBuilder().Secret(endToEndSecret).Issuer("auth").TTL(time.Hour).Build()
	require.NoError(t, err)

	tokenUseCase, err := tokenusecase.NewBuilder().
		TokenService(tokenService).
		LogPublisher(logPublisher).
		Build()
	require.NoError(t, err)

	controller, err := authgrpccontroller.NewBuilder().TokenUseCase(tokenUseCase).Build()
	require.NoError(t, err)

	server, err := grpcserver.NewBuilder().LogPublisher(logPublisher).Port(port).Build()
	require.NoError(t, err)

	authproto.RegisterAuthServiceServer(server.Server(), controller)

	serverCtx, stopServer := context.WithCancel(t.Context())
	runFailed := make(chan error, 1)

	go func() { runFailed <- server.Run(serverCtx) }()

	waitForPort(t, port)

	connection, err := grpc.NewClient(
		net.JoinHostPort("localhost", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = connection.Close() })

	gateway, err := commonauthgateway.NewBuilder().
		Client(authproto.NewAuthServiceClient(connection)).
		LogPublisher(logPublisher).
		Build()
	require.NoError(t, err)

	t.Run("настоящий токен проходит сквозь весь стек", func(t *testing.T) {
		expected := commonuserentity.NewUserID()

		issued, issueErr := tokenUseCase.Issue(t.Context(), expected)
		require.NoError(t, issueErr)
		require.NotEmpty(t, issued.AccessToken)

		userID, validateErr := gateway.ValidateToken(t.Context(), issued.AccessToken)

		require.NoError(t, validateErr)
		assert.Equal(t, expected, userID)
	})

	t.Run("мусор вместо токена отвергается", func(t *testing.T) {
		userID, validateErr := gateway.ValidateToken(t.Context(), "not-a-jwt")

		require.Error(t, validateErr)
		assert.True(t, userID.IsZero())
	})

	t.Run("токен чужого издателя отвергается", func(t *testing.T) {
		foreign, buildErr := jwt.NewBuilder().Secret(endToEndSecret).Issuer("blog").TTL(time.Hour).Build()
		require.NoError(t, buildErr)

		rawToken, _, issueErr := foreign.Issue(commonuserentity.NewUserID().String())
		require.NoError(t, issueErr)

		userID, validateErr := gateway.ValidateToken(t.Context(), rawToken)

		require.Error(t, validateErr)
		assert.True(t, userID.IsZero())
	})

	stopServer()
	require.NoError(t, <-runFailed)
}
