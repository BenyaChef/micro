package initservice

import "github.com/BenyaChef/micro/infrastructure/grpcserver"

func (c *DependencyContainer) initGRPCServer() error {
	server, err := grpcserver.NewBuilder().
		LogPublisher(c.infra.LogPublisher()).
		GRPCMonitor(c.infra.GRPCMonitor()).
		Port(c.environments.grpcPort).
		Build()
	if err != nil {
		return err
	}

	c.grpcServer = server

	return nil
}
