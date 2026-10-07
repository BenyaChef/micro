package gatewayinterface

import (
	"context"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

type AuthGateway interface {
	ValidateToken(ctx context.Context, rawToken string) (commonuserentity.UserID, error)
}
