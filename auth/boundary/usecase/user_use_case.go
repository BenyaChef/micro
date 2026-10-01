package usecaseinterface

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

type UserUseCase interface {
	Register(ctx context.Context, data *boundarymodel.RegisterData) (*userentity.User, error)
	Login(ctx context.Context, data *boundarymodel.LoginData) (*boundarymodel.LoginResult, error)
	GetByID(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error)
}
