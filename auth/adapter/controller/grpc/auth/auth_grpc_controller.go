package authgrpccontroller

import (
	"context"

	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	"github.com/BenyaChef/micro/infrastructure/grpcserver/grpcmodel"
)

var _ authproto.AuthServiceServer = (*AuthController)(nil)

type AuthController struct {
	authproto.UnimplementedAuthServiceServer

	tokenUseCase usecaseinterface.TokenUseCase
}

func (c *AuthController) ValidateToken(
	ctx context.Context,
	request *authproto.ValidateTokenRequest,
) (*authproto.ValidateTokenResponse, error) {
	userID, err := c.tokenUseCase.Validate(ctx, request.GetAccessToken())
	if err != nil {
		return errorResponse(err), nil
	}

	return &authproto.ValidateTokenResponse{
		Result: &authproto.ValidateTokenResponse_Token{
			Token: &authproto.ValidatedToken{UserId: userID.String()},
		},
	}, nil
}

func errorResponse(err error) *authproto.ValidateTokenResponse {
	return &authproto.ValidateTokenResponse{
		Result: &authproto.ValidateTokenResponse_Error{Error: grpcmodel.ToErrorModel(err)},
	}
}
