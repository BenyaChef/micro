package initservice

import (
	authgateway "github.com/BenyaChef/micro/blog/adapter/gateway/auth"
	commonauthgateway "github.com/BenyaChef/micro/common/auth/adapter/gateway"
	authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"
	"github.com/BenyaChef/micro/infrastructure/grpcclient"
)

func (c *DependencyContainer) initGateways() error {
	client, err := grpcclient.NewBuilder().
		Target(c.environments.authGRPCAddr).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.authGRPCClient = client

	commonGateway, err := commonauthgateway.NewBuilder().
		Client(authproto.NewAuthServiceClient(client.Connection())).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	gateway, err := authgateway.NewBuilder().
		AuthGateway(commonGateway).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.authGateway = gateway

	tokenValidator, err := authgateway.NewTokenValidator(gateway)
	if err != nil {
		return err
	}

	c.tokenValidator = tokenValidator

	return nil
}
