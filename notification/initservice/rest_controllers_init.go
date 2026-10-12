package initservice

import "github.com/BenyaChef/micro/infrastructure/restserver"

func (c *DependencyContainer) initRESTServer() error {
	server, err := restserver.NewBuilder().
		LogPublisher(c.infra.LogPublisher()).
		ServiceName(c.serviceName).
		Port(c.environments.restPort).
		Build()
	if err != nil {
		return err
	}

	c.restServer = server

	return nil
}
