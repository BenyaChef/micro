package initservice

import (
	"context"

	authrestcontroller "github.com/BenyaChef/micro/auth/adapter/controller/rest/auth"
	repositoryinterface "github.com/BenyaChef/micro/auth/boundary/repository"
	serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
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
	restServer     *restserver.Server

	tokenService   serviceinterface.TokenService
	passwordHasher serviceinterface.PasswordHasher

	userRepository repositoryinterface.UserRepository
	tokenUseCase   usecaseinterface.TokenUseCase
	userUseCase    usecaseinterface.UserUseCase
	authController *authrestcontroller.AuthController
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
		container.initServices,
		container.initUseCases,
		container.initRESTServer,
		container.initRESTControllers,
		container.initRESTRouting,
	}

	for _, initStep := range initChain {
		if err := initStep(); err != nil {
			container.Close()

			return nil, err
		}
	}

	return container, nil
}

func (c *DependencyContainer) RESTServer() restserverinterface.RESTServer {
	return c.restServer
}

func (c *DependencyContainer) Close() {
	if c.postgresClient != nil {
		c.postgresClient.Close()
	}
}
