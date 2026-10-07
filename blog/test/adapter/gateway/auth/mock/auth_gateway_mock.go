package authgatewaymock

import (
	"context"

	commonauthgatewayinterface "github.com/BenyaChef/micro/common/auth/boundary/gateway"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

var _ commonauthgatewayinterface.AuthGateway = (*AuthGatewayMock)(nil)

type AuthGatewayMock struct {
	ValidateTokenFunc func(ctx context.Context, rawToken string) (commonuserentity.UserID, error)
}

func NewAuthGatewayMock() *AuthGatewayMock {
	return &AuthGatewayMock{
		ValidateTokenFunc: func(context.Context, string) (commonuserentity.UserID, error) {
			return commonuserentity.NewUserID(), nil
		},
	}
}

func (m *AuthGatewayMock) ValidateToken(ctx context.Context, rawToken string) (commonuserentity.UserID, error) {
	return m.ValidateTokenFunc(ctx, rawToken)
}
