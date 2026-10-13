package otel

import (
	"context"
	"sync"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Provider struct {
	meterProvider *sdkmetric.MeterProvider
	metricsServer *MetricsServer
	logPublisher  loggerinterface.LogPublisher
	shutdownOnce  sync.Once
}

func (p *Provider) MeterProvider() *sdkmetric.MeterProvider {
	return p.meterProvider
}

func (p *Provider) MetricsServer() *MetricsServer {
	return p.metricsServer
}

func (p *Provider) Close(ctx context.Context) {
	if p == nil {
		return
	}

	p.shutdownOnce.Do(func() {
		if p.meterProvider == nil {
			return
		}

		if err := p.meterProvider.Shutdown(ctx); err != nil {
			p.logPublisher.LogError(ctx, err)
		}
	})
}
