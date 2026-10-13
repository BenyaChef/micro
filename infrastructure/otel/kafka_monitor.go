package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const kafkaMeterName = "micro/infrastructure/queueserver"

const (
	ResultSuccess = "success"
	ResultError   = "error"
)

type KafkaMonitor struct {
	publishCount    metric.Int64Counter
	processDuration metric.Float64Histogram
}

func NewKafkaMonitor(meterProvider *sdkmetric.MeterProvider) (*KafkaMonitor, error) {
	meter := meterProvider.Meter(kafkaMeterName)

	publishCount, err := meter.Int64Counter("messaging.publish.count",
		metric.WithDescription("Количество опубликованных сообщений"),
	)
	if err != nil {
		return nil, ErrMonitorFailed(kafkaMeterName, err)
	}

	processDuration, err := meter.Float64Histogram("messaging.process.duration",
		metric.WithUnit("s"),
		metric.WithDescription("Длительность обработки сообщения из очереди"),
	)
	if err != nil {
		return nil, ErrMonitorFailed(kafkaMeterName, err)
	}

	return &KafkaMonitor{publishCount: publishCount, processDuration: processDuration}, nil
}

func (m *KafkaMonitor) RecordPublish(ctx context.Context, topic, result string) {
	if m == nil {
		return
	}

	m.publishCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String("messaging.destination.name", topic),
		attribute.String("messaging.operation.result", result),
	))
}

func (m *KafkaMonitor) RecordProcess(ctx context.Context, topic, result string, seconds float64) {
	if m == nil {
		return
	}

	m.processDuration.Record(ctx, seconds, metric.WithAttributes(
		attribute.String("messaging.destination.name", topic),
		attribute.String("messaging.operation.result", result),
	))
}
