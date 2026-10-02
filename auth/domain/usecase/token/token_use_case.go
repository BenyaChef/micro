package tokenusecase

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

var _ usecaseinterface.TokenUseCase = (*TokenUseCase)(nil)

type TokenUseCase struct {
	tokenService   serviceinterface.TokenService
	logPublisher   loggerinterface.LogPublisher
	errorProcessor *errorProcessor
}

func (uc *TokenUseCase) Issue(ctx context.Context, userID commonuserentity.UserID) (*boundarymodel.LoginResult, error) {
	if userID.IsZero() {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrUserIDIsRequired)
	}

	rawToken, expiresIn, err := uc.tokenService.Issue(userID.String())
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturnIssueError(ctx, err)
	}

	return &boundarymodel.LoginResult{
		AccessToken: rawToken,
		ExpiresIn:   int(expiresIn.Seconds()),
	}, nil
}

func (uc *TokenUseCase) Validate(ctx context.Context, rawToken string) (commonuserentity.UserID, error) {
	if rawToken == "" {
		return commonuserentity.UserID{}, uc.errorProcessor.LogAndReturn(ctx, ErrTokenIsRequired)
	}

	subject, err := uc.tokenService.Subject(rawToken)
	if err != nil {
		return commonuserentity.UserID{}, uc.errorProcessor.LogAndReturn(ctx, err)
	}

	userID, err := commonuserentity.UserIDFrom(subject)
	if err != nil {
		return commonuserentity.UserID{}, uc.errorProcessor.LogAndReturn(ctx, err)
	}

	return userID, nil
}
