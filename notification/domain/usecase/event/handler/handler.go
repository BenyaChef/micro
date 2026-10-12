package eventusecasehandler

import (
	"context"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
)

type Handler interface {
	Handle(ctx context.Context, event *commonautheventmodel.UserEvent) error
}
