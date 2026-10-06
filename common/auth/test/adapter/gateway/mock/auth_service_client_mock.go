package commonauthgatewaymock

import (
	"context"

	"google.golang.org/grpc"

	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
)

var _ authproto.AuthServiceClient = (*AuthServiceClientMock)(nil)

type AuthServiceClientMock struct {
	ValidateTokenFunc func(ctx context.Context, request *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error)
}

func NewAuthServiceClientMock() *AuthServiceClientMock {
	return &AuthServiceClientMock{
		ValidateTokenFunc: func(context.Context, *authproto.ValidateTokenRequest) (*authproto.ValidateTokenResponse, error) {
			return &authproto.ValidateTokenResponse{}, nil
		},
	}
}

func (m *AuthServiceClientMock) ValidateToken(
	ctx context.Context,
	request *authproto.ValidateTokenRequest,
	_ ...grpc.CallOption,
) (*authproto.ValidateTokenResponse, error) {
	return m.ValidateTokenFunc(ctx, request)
}
