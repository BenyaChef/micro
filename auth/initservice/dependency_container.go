package initservice

import (
	"context"

	authgrpccontroller "github.com/BenyaChef/micro/auth/adapter/controller/grpc/auth"
	authrestcontroller "github.com/BenyaChef/micro/auth/adapter/controller/rest/auth"
	eventinterface "github.com/BenyaChef/micro/auth/boundary/event"
	repositoryinterface "github.com/BenyaChef/micro/auth/boundary/repository"
	serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	"github.com/BenyaChef/micro/infrastructure/grpcserver"
	grpcserverinterface "github.com/BenyaChef/micro/infrastructure/grpcserver/interface"
	"github.com/BenyaChef/micro/infrastructure/initinfra"
	"github.com/BenyaChef/micro/infrastructure/postgres"
	"github.com/BenyaChef/micro/infrastructure/queueserver"
	"github.com/BenyaChef/micro/infrastructure/restserver"
	restserverinterface "github.com/BenyaChef/micro/infrastructure/restserver/interface"
)

type DependencyContainer struct {
	infra       *initinfra.Container
	serviceName string

	environments *environments

	postgresClient *postgres.Client
	queuePublisher *queueserver.Publisher
	restServer     *restserver.Server
	grpcServer     *grpcserver.Server

	tokenService   serviceinterface.TokenService
	passwordHasher serviceinterface.PasswordHasher

	userRepository     repositoryinterface.UserRepository
	userEventPublisher eventinterface.UserEventPublisher
	tokenUseCase       usecaseinterface.TokenUseCase
	userUseCase        usecaseinterface.UserUseCase
	authController     *authrestcontroller.AuthController

	authGRPCController *authgrpccontroller.AuthController
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
		container.initQueuePublisher,
		container.initEventPublishers,
		container.initServices,
		container.initUseCases,
		container.initRESTServer,
		container.initRESTControllers,
		container.initRESTRouting,
		container.initGRPCServer,
		container.initGRPCControllers,
		container.initGRPCRouting,
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

func (c *DependencyContainer) GRPCServer() grpcserverinterface.GRPCServer {
	return c.grpcServer
}

func (c *DependencyContainer) Close(ctx context.Context) {
	if c.queuePublisher != nil {
		c.queuePublisher.Close(ctx)
	}

	if c.postgresClient != nil {
		c.postgresClient.Close()
	}
}
