package postgresmodel

import "time"

const (
	DefaultMaxConns int32 = 10
	DefaultMinConns int32 = 2

	DefaultConnectTimeout    = 5 * time.Second
	DefaultMaxConnLifetime   = 30 * time.Minute
	DefaultMaxConnIdleTime   = 5 * time.Minute
	DefaultHealthCheckPeriod = time.Minute
)
