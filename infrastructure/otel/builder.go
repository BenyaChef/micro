package otel

import (
	"context"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const DefaultMetricsPort = "9101"

type Builder struct {
	serviceName  string
	environment  string
	metricsPort  string
	logPublisher loggerinterface.LogPublisher
}

func NewBuilder() *Builder {
	return &Builder{metricsPort: DefaultMetricsPort}
}

func (b *Builder) ServiceName(serviceName string) *Builder {
	b.serviceName = serviceName

	return b
}

func (b *Builder) Environment(environment string) *Builder {
	b.environment = environment

	return b
}

func (b *Builder) MetricsPort(metricsPort string) *Builder {
	b.metricsPort = metricsPort

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) Build(ctx context.Context) (*Provider, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	otelResource, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(b.serviceName),
		semconv.DeploymentEnvironmentNameKey.String(b.environment),
	))
	if err != nil {
		return nil, ErrResourceFailed(err)
	}

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	exporter, err := promexporter.New(
		promexporter.WithRegisterer(promRegistry),
		promexporter.WithoutScopeInfo(),
	)
	if err != nil {
		return nil, ErrExporterFailed(err)
	}

	options := []sdkmetric.Option{
		sdkmetric.WithResource(otelResource),
		sdkmetric.WithReader(exporter),
	}

	for _, view := range AllViews() {
		options = append(options, sdkmetric.WithView(view))
	}

	meterProvider := sdkmetric.NewMeterProvider(options...)
	handler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})

	return &Provider{
		meterProvider: meterProvider,
		metricsServer: NewMetricsServer(handler, b.metricsPort, b.logPublisher),
		logPublisher:  b.logPublisher,
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.serviceName == "" {
		errs = append(errs, ErrServiceNameIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if b.metricsPort == "" {
		errs = append(errs, ErrPortIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
