package initinfra

import (
	"context"

	"github.com/BenyaChef/micro/infrastructure/envregistry"
	envregistryinterface "github.com/BenyaChef/micro/infrastructure/envregistry/interface"
	envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"
	"github.com/BenyaChef/micro/infrastructure/logger"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	infraotel "github.com/BenyaChef/micro/infrastructure/otel"
)

const (
	defaultLogLevel    = "info"
	defaultEnvironment = "local"
)

type Container struct {
	logPublisher loggerinterface.LogPublisher
	env          envregistryinterface.EnvRegistry

	otelProvider *infraotel.Provider
	httpMonitor  *infraotel.HTTPMonitor
	grpcMonitor  *infraotel.GRPCMonitor
	kafkaMonitor *infraotel.KafkaMonitor
}

func Bootstrap(ctx context.Context, serviceName string) (*Container, error) {
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

	container := &Container{logPublisher: logPublisher, env: env}

	if err := container.initOTel(ctx, serviceName); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initOTel(ctx context.Context, serviceName string) error {
	provider, err := infraotel.NewBuilder().
		ServiceName(serviceName).
		Environment(c.env.StringDefault(envregistrymodel.KeyEnvironment, defaultEnvironment)).
		MetricsPort(c.env.StringDefault(envregistrymodel.KeyMetricsPort, infraotel.DefaultMetricsPort)).
		LogPublisher(c.logPublisher).
		Build(ctx)
	if err != nil {
		return err
	}

	c.otelProvider = provider

	httpMonitor, err := infraotel.NewHTTPMonitor(provider.MeterProvider())
	if err != nil {
		return err
	}

	grpcMonitor, err := infraotel.NewGRPCMonitor(provider.MeterProvider())
	if err != nil {
		return err
	}

	kafkaMonitor, err := infraotel.NewKafkaMonitor(provider.MeterProvider())
	if err != nil {
		return err
	}

	c.httpMonitor = httpMonitor
	c.grpcMonitor = grpcMonitor
	c.kafkaMonitor = kafkaMonitor

	return nil
}

func (c *Container) LogPublisher() loggerinterface.LogPublisher {
	return c.logPublisher
}

func (c *Container) Env() envregistryinterface.EnvRegistry {
	return c.env
}

func (c *Container) OTelProvider() *infraotel.Provider {
	return c.otelProvider
}

func (c *Container) MetricsServer() *infraotel.MetricsServer {
	return c.otelProvider.MetricsServer()
}

func (c *Container) HTTPMonitor() *infraotel.HTTPMonitor {
	return c.httpMonitor
}

func (c *Container) GRPCMonitor() *infraotel.GRPCMonitor {
	return c.grpcMonitor
}

func (c *Container) KafkaMonitor() *infraotel.KafkaMonitor {
	return c.kafkaMonitor
}

func (c *Container) Close(ctx context.Context) {
	if c.otelProvider != nil {
		c.otelProvider.Close(ctx)
	}
}
