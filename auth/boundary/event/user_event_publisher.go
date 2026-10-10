package eventinterface

import (
	"context"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
)

type UserEventPublisher interface {
	Registered(ctx context.Context, userID commonuserentity.UserID, email emailprimitive.Email) error
}
