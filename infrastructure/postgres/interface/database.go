package postgresinterface

import "context"

type Database interface {
	Executor

	Ping(ctx context.Context) error
	Close()
}
