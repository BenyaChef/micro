package postrestcontrollermock

import (
	"context"

	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
)

var _ restserverinterface.TokenValidator = (*TokenValidatorMock)(nil)

type TokenValidatorMock struct {
	ValidateTokenFunc func(ctx context.Context, rawToken string) (string, error)
}

func NewTokenValidatorMock(callerID string) *TokenValidatorMock {
	return &TokenValidatorMock{
		ValidateTokenFunc: func(context.Context, string) (string, error) { return callerID, nil },
	}
}

func (m *TokenValidatorMock) ValidateToken(ctx context.Context, rawToken string) (string, error) {
	return m.ValidateTokenFunc(ctx, rawToken)
}
