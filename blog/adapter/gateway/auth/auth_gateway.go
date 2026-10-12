package authgateway

import (
	"context"

	gatewayinterface "github.com/BenyaChef/micro/blog/boundary/gateway"
	commonauthgateway "github.com/BenyaChef/micro/common/auth/adapter/gateway"
	commonauthgatewayinterface "github.com/BenyaChef/micro/common/auth/boundary/gateway"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

var _ gatewayinterface.AuthGateway = (*AuthGateway)(nil)

type AuthGateway struct {
	authGateway  commonauthgatewayinterface.AuthGateway
	logPublisher loggerinterface.LogPublisher
}

func (g *AuthGateway) ValidateToken(ctx context.Context, rawToken string) (commonuserentity.UserID, error) {
	userID, err := g.authGateway.ValidateToken(ctx, rawToken)
	if err == nil {
		return userID, nil
	}

	if errors.EqualByCode(err, commonauthgateway.ErrCodeCallFailed) {
		return commonuserentity.UserID{}, ErrAuthUnavailable(err)
	}

	return commonuserentity.UserID{}, ErrTokenRejected(err)
}
