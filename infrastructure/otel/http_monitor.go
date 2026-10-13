package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const httpMeterName = "micro/infrastructure/restserver"

type HTTPMonitor struct {
	requestDuration metric.Float64Histogram
}

func NewHTTPMonitor(meterProvider *sdkmetric.MeterProvider) (*HTTPMonitor, error) {
	meter := meterProvider.Meter(httpMeterName)

	requestDuration, err := meter.Float64Histogram("http.server.request.duration",
		metric.WithUnit("s"),
		metric.WithDescription("Длительность обработки HTTP-запроса"),
	)
	if err != nil {
		return nil, ErrMonitorFailed(httpMeterName, err)
	}

	return &HTTPMonitor{requestDuration: requestDuration}, nil
}

func (m *HTTPMonitor) RecordRequest(ctx context.Context, method, route string, status int, seconds float64) {
	if m == nil {
		return
	}

	m.requestDuration.Record(ctx, seconds, metric.WithAttributes(
		attribute.String("http.request.method", method),
		attribute.String("http.route", route),
		attribute.Int("http.response.status_code", status),
	))
}
