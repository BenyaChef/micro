package initservice

import (
	"context"

	postrepository "github.com/BenyaChef/micro/blog/adapter/repository/post"
	infraotel "github.com/BenyaChef/micro/infrastructure/otel"
	"github.com/BenyaChef/micro/infrastructure/postgres"
)

func (c *DependencyContainer) initPostgres(ctx context.Context) error {
	builder := postgres.NewBuilder().DSN(c.environments.postgresDSN)

	if c.environments.postgresMaxConns > 0 {
		builder = builder.MaxConns(int32(c.environments.postgresMaxConns))
	}

	if c.environments.postgresMinConns > 0 {
		builder = builder.MinConns(int32(c.environments.postgresMinConns))
	}

	client, err := builder.Build(ctx)
	if err != nil {
		return err
	}

	c.postgresClient = client

	return infraotel.RegisterPostgresMonitor(c.infra.OTelProvider().MeterProvider(), client.Pool())
}

func (c *DependencyContainer) initRepositories() error {
	postRepository, err := postrepository.NewBuilder().
		Executor(c.postgresClient).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.postRepository = postRepository

	return nil
}
