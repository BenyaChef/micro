package authgateway

import (
	"context"

	gatewayinterface "github.com/BenyaChef/micro/blog/boundary/gateway"
	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
)

var _ restserverinterface.TokenValidator = (*TokenValidator)(nil)

type TokenValidator struct {
	authGateway gatewayinterface.AuthGateway
}

func NewTokenValidator(authGateway gatewayinterface.AuthGateway) (*TokenValidator, error) {
	if authGateway == nil {
		return nil, ErrAuthGatewayIsRequired
	}

	return &TokenValidator{authGateway: authGateway}, nil
}

func (v *TokenValidator) ValidateToken(ctx context.Context, rawToken string) (string, error) {
	userID, err := v.authGateway.ValidateToken(ctx, rawToken)
	if err != nil {
		return "", err
	}

	return userID.String(), nil
}
