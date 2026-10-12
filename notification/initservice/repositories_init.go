package initservice

import (
	"context"

	"github.com/BenyaChef/micro/infrastructure/postgres"
	notificationrepository "github.com/BenyaChef/micro/notification/adapter/repository/notification"
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

	return nil
}

func (c *DependencyContainer) initRepositories() error {
	repository, err := notificationrepository.NewBuilder().
		Executor(c.postgresClient).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.notificationRepository = repository

	return nil
}
