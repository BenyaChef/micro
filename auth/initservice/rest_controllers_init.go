package initservice

import (
	authrestcontroller "github.com/BenyaChef/micro/auth/adapter/controller/rest/auth"
	"github.com/BenyaChef/micro/infrastructure/restserver"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
)

func (c *DependencyContainer) initRESTServer() error {
	server, err := restserver.NewBuilder().
		LogPublisher(c.infra.LogPublisher()).
		ErrorResolver(authrestcontroller.NewErrorResolver()).
		ServiceName(c.serviceName).
		Port(c.environments.restPort).
		Build()
	if err != nil {
		return err
	}

	c.restServer = server

	return nil
}

func (c *DependencyContainer) initRESTControllers() error {
	baseController, err := restservercontroller.NewBaseController(
		c.restServer.ResponseService(), c.infra.LogPublisher(),
	)
	if err != nil {
		return err
	}

	authController, err := authrestcontroller.NewBuilder().
		BaseController(baseController).
		UserUseCase(c.userUseCase).
		TokenUseCase(c.tokenUseCase).
		Build()
	if err != nil {
		return err
	}

	c.authController = authController

	return nil
}
