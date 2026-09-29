package initinfra

import (
	"micro/infrastructure/envregistry"
	envregistryinterface "micro/infrastructure/envregistry/interface"
	envregistrymodel "micro/infrastructure/envregistry/model"
	"micro/infrastructure/logger"
	loggerinterface "micro/infrastructure/logger/interface"
)

type Container struct {
	logPublisher loggerinterface.LogPublisher
	env          envregistryinterface.EnvRegistry
}

func Bootstrap(serviceName string) (*Container, error) {
	if serviceName == "" {
		return nil, ErrServiceNameIsRequired
	}

	env, err := envregistry.NewBuilder().Build()
	if err != nil {
		return nil, err
	}

	logPublisher, err := logger.NewBuilder().
		ServiceName(serviceName).
		Level(env.StringDefault(envregistrymodel.KeyLogLevel, defaultLogLevel)).
		Build()
	if err != nil {
		return nil, err
	}

	return &Container{logPublisher: logPublisher, env: env}, nil
}

func (c *Container) LogPublisher() loggerinterface.LogPublisher {
	return c.logPublisher
}

func (c *Container) Env() envregistryinterface.EnvRegistry {
	return c.env
}
