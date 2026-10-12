package initservice

import envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"

const (
	defaultRESTPort     = "8083"
	defaultKafkaBrokers = "localhost:29092"
	defaultKafkaGroupID = "notification-service"
)

type environments struct {
	restPort         string
	postgresDSN      string
	postgresMaxConns int
	postgresMinConns int
	kafkaBrokers     string
	kafkaGroupID     string
}

func (c *DependencyContainer) initEnvironments() error {
	env := c.infra.Env()

	c.environments = &environments{
		restPort:         env.StringDefault(envregistrymodel.KeyRestPort, defaultRESTPort),
		postgresDSN:      env.String(envregistrymodel.KeyPostgresDSN),
		postgresMaxConns: env.Int(envregistrymodel.KeyPostgresMaxConns, 0),
		postgresMinConns: env.Int(envregistrymodel.KeyPostgresMinConns, 0),
		kafkaBrokers:     env.StringDefault(envregistrymodel.KeyKafkaBrokers, defaultKafkaBrokers),
		kafkaGroupID:     env.StringDefault(envregistrymodel.KeyKafkaGroupID, defaultKafkaGroupID),
	}

	return env.Err()
}
