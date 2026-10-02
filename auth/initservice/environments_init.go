package initservice

import (
	"time"

	envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"
	jwtmodel "github.com/BenyaChef/micro/infrastructure/jwt/model"
)

const (
	defaultRESTPort = "8081"
)

type environments struct {
	restPort         string
	postgresDSN      string
	postgresMaxConns int
	postgresMinConns int
	jwtSecret        string
	jwtTTL           time.Duration
	passwordHashCost int
}

func (c *DependencyContainer) initEnvironments() error {
	env := c.infra.Env()

	c.environments = &environments{
		restPort:         env.StringDefault(envregistrymodel.KeyRestPort, defaultRESTPort),
		postgresDSN:      env.String(envregistrymodel.KeyPostgresDSN),
		postgresMaxConns: env.Int(envregistrymodel.KeyPostgresMaxConns, 0),
		postgresMinConns: env.Int(envregistrymodel.KeyPostgresMinConns, 0),
		jwtSecret:        env.String(envregistrymodel.KeyJWTSecret),
		jwtTTL:           env.Duration(envregistrymodel.KeyJWTTTL, jwtmodel.DefaultTTL),
		passwordHashCost: env.Int(envregistrymodel.KeyPasswordHashCost, 0),
	}

	return env.Err()
}
