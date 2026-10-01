package usecaseinterface

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

type TokenUseCase interface {
	Issue(ctx context.Context, userID commonuserentity.UserID) (*boundarymodel.LoginResult, error)
	Validate(ctx context.Context, rawToken string) (commonuserentity.UserID, error)
}
