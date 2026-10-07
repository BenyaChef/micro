package initservice

import envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"

const (
	defaultRESTPort     = "8082"
	defaultAuthGRPCAddr = "localhost:9091"
)

type environments struct {
	restPort         string
	postgresDSN      string
	postgresMaxConns int
	postgresMinConns int
	authGRPCAddr     string
}

func (c *DependencyContainer) initEnvironments() error {
	env := c.infra.Env()

	c.environments = &environments{
		restPort:         env.StringDefault(envregistrymodel.KeyRestPort, defaultRESTPort),
		postgresDSN:      env.String(envregistrymodel.KeyPostgresDSN),
		postgresMaxConns: env.Int(envregistrymodel.KeyPostgresMaxConns, 0),
		postgresMinConns: env.Int(envregistrymodel.KeyPostgresMinConns, 0),
		authGRPCAddr:     env.StringDefault(envregistrymodel.KeyAuthGRPCAddr, defaultAuthGRPCAddr),
	}

	return env.Err()
}
