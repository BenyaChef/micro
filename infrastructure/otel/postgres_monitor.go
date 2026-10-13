package otel

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const postgresMeterName = "micro/infrastructure/postgres"

// Регистрирует наблюдаемый gauge: значения снимаются в момент скрейпа,
// отдельный горутин-семплер не нужен. Объекта, который кто-то дальше
// вызывает, тут нет, поэтому возвращаем только ошибку.
func RegisterPostgresMonitor(meterProvider *sdkmetric.MeterProvider, pool *pgxpool.Pool) error {
	meter := meterProvider.Meter(postgresMeterName)

	_, err := meter.Int64ObservableGauge("db.client.connection.count",
		metric.WithDescription("Соединения в пуле по состояниям"),
		metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
			if pool == nil {
				return nil
			}

			stat := pool.Stat()

			observer.Observe(int64(stat.AcquiredConns()), metric.WithAttributes(connectionState("used")))
			observer.Observe(int64(stat.IdleConns()), metric.WithAttributes(connectionState("idle")))
			observer.Observe(int64(stat.TotalConns()), metric.WithAttributes(connectionState("total")))
			observer.Observe(int64(stat.MaxConns()), metric.WithAttributes(connectionState("max")))

			return nil
		}),
	)
	if err != nil {
		return ErrMonitorFailed(postgresMeterName, err)
	}

	return nil
}

func connectionState(state string) attribute.KeyValue {
	return attribute.String("db.client.connection.state", state)
}
