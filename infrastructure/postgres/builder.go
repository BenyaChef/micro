package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	postgresmodel "micro/infrastructure/postgres/model"
)

type Builder struct {
	dsn               string
	maxConns          int32
	minConns          int32
	connectTimeout    time.Duration
	maxConnLifetime   time.Duration
	maxConnIdleTime   time.Duration
	healthCheckPeriod time.Duration
}

func NewBuilder() *Builder {
	return &Builder{
		maxConns:          postgresmodel.DefaultMaxConns,
		minConns:          postgresmodel.DefaultMinConns,
		connectTimeout:    postgresmodel.DefaultConnectTimeout,
		maxConnLifetime:   postgresmodel.DefaultMaxConnLifetime,
		maxConnIdleTime:   postgresmodel.DefaultMaxConnIdleTime,
		healthCheckPeriod: postgresmodel.DefaultHealthCheckPeriod,
	}
}

func (b *Builder) DSN(dsn string) *Builder {
	b.dsn = dsn

	return b
}

func (b *Builder) MaxConns(maxConns int32) *Builder {
	b.maxConns = maxConns

	return b
}

func (b *Builder) MinConns(minConns int32) *Builder {
	b.minConns = minConns

	return b
}

func (b *Builder) ConnectTimeout(timeout time.Duration) *Builder {
	b.connectTimeout = timeout

	return b
}

func (b *Builder) MaxConnLifetime(lifetime time.Duration) *Builder {
	b.maxConnLifetime = lifetime

	return b
}

func (b *Builder) MaxConnIdleTime(idleTime time.Duration) *Builder {
	b.maxConnIdleTime = idleTime

	return b
}

func (b *Builder) HealthCheckPeriod(period time.Duration) *Builder {
	b.healthCheckPeriod = period

	return b
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.dsn == "" {
		errs = append(errs, ErrDSNIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (b *Builder) Build(ctx context.Context) (*Client, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	config, err := pgxpool.ParseConfig(b.dsn)
	if err != nil {
		return nil, ErrInvalidDSN(err)
	}

	config.MaxConns = b.maxConns
	config.MinConns = b.minConns
	config.MaxConnLifetime = b.maxConnLifetime
	config.MaxConnIdleTime = b.maxConnIdleTime
	config.HealthCheckPeriod = b.healthCheckPeriod
	config.ConnConfig.ConnectTimeout = b.connectTimeout

	host := config.ConnConfig.Host

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, ErrConnectFailed(host, err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, b.connectTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()

		return nil, ErrPingFailed(host, err)
	}

	return &Client{pool: pool, host: host}, nil
}
