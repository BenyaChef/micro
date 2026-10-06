package initservice

import authproto "github.com/BenyaChef/micro/common/auth/adapter/gateway/proto"

func (c *DependencyContainer) initGRPCRouting() error {
	authproto.RegisterAuthServiceServer(c.grpcServer.Server(), c.authGRPCController)

	return nil
}
