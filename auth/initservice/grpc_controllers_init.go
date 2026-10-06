package initservice

import authgrpccontroller "github.com/BenyaChef/micro/auth/adapter/controller/grpc/auth"

func (c *DependencyContainer) initGRPCControllers() error {
	controller, err := authgrpccontroller.NewBuilder().
		TokenUseCase(c.tokenUseCase).
		Build()
	if err != nil {
		return err
	}

	c.authGRPCController = controller

	return nil
}
