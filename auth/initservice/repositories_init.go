package initservice

import (
	"context"

	userrepository "github.com/BenyaChef/micro/auth/adapter/repository/user"
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

	return nil
}

func (c *DependencyContainer) initRepositories() error {
	userRepository, err := userrepository.NewBuilder().
		Executor(c.postgresClient).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.userRepository = userRepository

	return nil
}
