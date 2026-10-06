package commonauthgateway

import (
	"context"

	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	commonauthgatewayinterface "github.com/BenyaChef/micro/common/auth/boundary/gateway"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	"github.com/BenyaChef/micro/infrastructure/grpcserver/grpcmodel"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

var _ commonauthgatewayinterface.AuthGateway = (*AuthGRPCGateway)(nil)

type AuthGRPCGateway struct {
	client         authproto.AuthServiceClient
	logPublisher   loggerinterface.LogPublisher
	errorProcessor *errorProcessor
}

func (g *AuthGRPCGateway) ValidateToken(ctx context.Context, rawToken string) (commonuserentity.UserID, error) {
	if rawToken == "" {
		return commonuserentity.UserID{}, g.errorProcessor.LogAndReturn(ctx, ErrTokenIsRequired)
	}

	request := &authproto.ValidateTokenRequest{AccessToken: rawToken}

	response, err := g.client.ValidateToken(ctx, request)
	if err != nil {
		return commonuserentity.UserID{}, g.errorProcessor.LogAndReturnCallError(ctx, err)
	}

	if errModel := response.GetError(); errModel != nil {
		return commonuserentity.UserID{}, g.errorProcessor.LogAndReturn(ctx, grpcmodel.FromErrorModel(errModel))
	}

	token := response.GetToken()
	if token == nil {
		return commonuserentity.UserID{}, g.errorProcessor.LogAndReturn(ctx, ErrEmptyResult)
	}

	userID, err := commonuserentity.UserIDFrom(token.GetUserId())
	if err != nil {
		return commonuserentity.UserID{}, g.errorProcessor.LogAndReturn(ctx, err)
	}

	return userID, nil
}
