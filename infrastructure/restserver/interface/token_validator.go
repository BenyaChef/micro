package restserverinterface

import "context"

type TokenValidator interface {
	ValidateToken(ctx context.Context, rawToken string) (callerID string, err error)
}
