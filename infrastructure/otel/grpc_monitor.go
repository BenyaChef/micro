package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const grpcMeterName = "micro/infrastructure/grpcserver"

type GRPCMonitor struct {
	requestDuration metric.Float64Histogram
}

func NewGRPCMonitor(meterProvider *sdkmetric.MeterProvider) (*GRPCMonitor, error) {
	meter := meterProvider.Meter(grpcMeterName)

	requestDuration, err := meter.Float64Histogram("rpc.server.duration",
		metric.WithUnit("s"),
		metric.WithDescription("Длительность обработки gRPC-вызова"),
	)
	if err != nil {
		return nil, ErrMonitorFailed(grpcMeterName, err)
	}

	return &GRPCMonitor{requestDuration: requestDuration}, nil
}

func (m *GRPCMonitor) RecordRequest(ctx context.Context, method, code string, seconds float64) {
	if m == nil {
		return
	}

	m.requestDuration.Record(ctx, seconds, metric.WithAttributes(
		attribute.String("rpc.method", method),
		attribute.String("rpc.grpc.status_code", code),
	))
}
