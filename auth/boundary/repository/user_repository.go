package repositoryinterface

import (
	"context"

	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
)

type UserRepository interface {
	Insert(ctx context.Context, user *userentity.User) error
	GetByID(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error)
	GetByEmail(ctx context.Context, email emailprimitive.Email) (*userentity.User, error)
}
