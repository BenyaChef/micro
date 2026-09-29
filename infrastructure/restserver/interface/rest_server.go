package restserverinterface

import (
	"context"

	"github.com/go-chi/chi/v5"
)

type RESTServer interface {
	Router() chi.Router
	Addr() string
	Run(ctx context.Context) error
}
