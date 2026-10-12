package usecaseinterface

import (
	"context"

	commonautheventmodel "github.com/BenyaChef/micro/common/auth/adapter/event/model"
)

type UserEventUseCase interface {
	Handle(ctx context.Context, event *commonautheventmodel.UserEvent) error
}
