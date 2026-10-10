package envregistrymodel

type EnvKey string

func (k EnvKey) String() string {
	return string(k)
}

const (
	KeyLogLevel EnvKey = "LOG_LEVEL"
	KeyRestPort EnvKey = "REST_PORT"
	KeyGRPCPort EnvKey = "GRPC_PORT"
)

const (
	KeyPostgresDSN            EnvKey = "POSTGRES_DSN"
	KeyPostgresMaxConns       EnvKey = "POSTGRES_MAX_CONNS"
	KeyPostgresMinConns       EnvKey = "POSTGRES_MIN_CONNS"
	KeyPostgresConnectTimeout EnvKey = "POSTGRES_CONNECT_TIMEOUT"
)

const (
	KeyJWTSecret EnvKey = "JWT_SECRET"
	KeyJWTTTL    EnvKey = "JWT_TTL"
)

const (
	KeyPasswordHashCost EnvKey = "PASSWORD_HASH_COST"
)

const (
	KeyAuthGRPCAddr EnvKey = "AUTH_GRPC_ADDR"
)

const (
	KeyKafkaBrokers EnvKey = "KAFKA_BROKERS"
)
