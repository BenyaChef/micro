package initservice

import (
	postrestcontroller "github.com/BenyaChef/micro/blog/adapter/controller/rest/post"
	"github.com/BenyaChef/micro/infrastructure/restserver"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
)

func (c *DependencyContainer) initRESTServer() error {
	server, err := restserver.NewBuilder().
		LogPublisher(c.infra.LogPublisher()).
		ErrorResolver(postrestcontroller.NewErrorResolver()).
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

	postController, err := postrestcontroller.NewBuilder().
		BaseController(baseController).
		PostUseCase(c.postUseCase).
		Build()
	if err != nil {
		return err
	}

	c.postController = postController

	return nil
}
