package restservermodel

import "time"

const (
	DefaultReadTimeout     = 10 * time.Second
	DefaultWriteTimeout    = 15 * time.Second
	DefaultIdleTimeout     = 60 * time.Second
	DefaultShutdownTimeout = 15 * time.Second
)
