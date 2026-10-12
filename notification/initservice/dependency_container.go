package initservice

import (
	"context"

	"github.com/BenyaChef/micro/infrastructure/initinfra"
	"github.com/BenyaChef/micro/infrastructure/postgres"
	"github.com/BenyaChef/micro/infrastructure/queueserver"
	"github.com/BenyaChef/micro/infrastructure/restserver"
	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
	repositoryinterface "github.com/BenyaChef/micro/notification/boundary/repository"
	usecaseinterface "github.com/BenyaChef/micro/notification/boundary/usecase"
)

type DependencyContainer struct {
	infra       *initinfra.Container
	serviceName string

	environments *environments

	postgresClient *postgres.Client
	queueConsumer  *queueserver.Consumer
	restServer     *restserver.Server

	notificationRepository repositoryinterface.NotificationRepository
	userEventUseCase       usecaseinterface.UserEventUseCase
}

func NewDependencyContainer(
	ctx context.Context,
	infra *initinfra.Container,
	serviceName string,
) (*DependencyContainer, error) {
	if infra == nil {
		return nil, ErrInfraContainerIsRequired
	}

	container := &DependencyContainer{infra: infra, serviceName: serviceName}

	initChain := []func() error{
		container.initEnvironments,
		func() error { return container.initPostgres(ctx) },
		container.initRepositories,
		container.initUseCases,
		container.initEventListeners,
		container.initRESTServer,
	}

	for _, initStep := range initChain {
		if err := initStep(); err != nil {
			container.Close(ctx)

			return nil, err
		}
	}

	return container, nil
}

func (c *DependencyContainer) RESTServer() restserverinterface.RESTServer {
	return c.restServer
}

func (c *DependencyContainer) QueueConsumer() *queueserver.Consumer {
	return c.queueConsumer
}

func (c *DependencyContainer) Close(ctx context.Context) {
	if c.queueConsumer != nil {
		c.queueConsumer.Close(ctx)
	}

	if c.postgresClient != nil {
		c.postgresClient.Close()
	}
}
