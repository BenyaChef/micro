package initservice

import (
	"context"

	postrestcontroller "github.com/BenyaChef/micro/blog/adapter/controller/rest/post"
	gatewayinterface "github.com/BenyaChef/micro/blog/boundary/gateway"
	repositoryinterface "github.com/BenyaChef/micro/blog/boundary/repository"
	usecaseinterface "github.com/BenyaChef/micro/blog/boundary/usecase"
	"github.com/BenyaChef/micro/infrastructure/grpcclient"
	"github.com/BenyaChef/micro/infrastructure/initinfra"
	"github.com/BenyaChef/micro/infrastructure/postgres"
	"github.com/BenyaChef/micro/infrastructure/restserver"
	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
)

type DependencyContainer struct {
	infra       *initinfra.Container
	serviceName string

	environments *environments

	postgresClient *postgres.Client
	authGRPCClient *grpcclient.Client
	restServer     *restserver.Server

	postRepository repositoryinterface.PostRepository
	authGateway    gatewayinterface.AuthGateway
	tokenValidator restserverinterface.TokenValidator
	postUseCase    usecaseinterface.PostUseCase
	postController *postrestcontroller.PostController
}

func NewDependencyContainer(ctx context.Context, infra *initinfra.Container, serviceName string) (*DependencyContainer, error) {
	if infra == nil {
		return nil, ErrInfraContainerIsRequired
	}

	container := &DependencyContainer{infra: infra, serviceName: serviceName}

	initChain := []func() error{
		container.initEnvironments,
		func() error { return container.initPostgres(ctx) },
		container.initRepositories,
		container.initGateways,
		container.initUseCases,
		container.initRESTServer,
		container.initRESTControllers,
		container.initRESTRouting,
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

func (c *DependencyContainer) Close(ctx context.Context) {
	if c.authGRPCClient != nil {
		c.authGRPCClient.Close(ctx)
	}

	if c.postgresClient != nil {
		c.postgresClient.Close()
	}
}
