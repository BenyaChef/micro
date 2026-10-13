package otel

import sdkmetric "go.opentelemetry.io/otel/sdk/metric"

var (
	httpDurationBuckets = []float64{
		0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
	}

	grpcDurationBuckets = []float64{
		0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.25, 0.5, 1, 5, 10,
	}

	kafkaDurationBuckets = []float64{
		0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 30,
	}

	postgresDurationBuckets = []float64{
		0.0005, 0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 5, 10,
	}
)

func HTTPViews() []sdkmetric.View {
	return []sdkmetric.View{bucketView("http.server.request.duration", httpDurationBuckets)}
}

func GRPCViews() []sdkmetric.View {
	return []sdkmetric.View{bucketView("rpc.server.duration", grpcDurationBuckets)}
}

func KafkaViews() []sdkmetric.View {
	return []sdkmetric.View{bucketView("messaging.*.duration", kafkaDurationBuckets)}
}

func PostgresViews() []sdkmetric.View {
	return []sdkmetric.View{bucketView("db.client.operation.duration", postgresDurationBuckets)}
}

func AllViews() []sdkmetric.View {
	var views []sdkmetric.View

	views = append(views, HTTPViews()...)
	views = append(views, GRPCViews()...)
	views = append(views, KafkaViews()...)
	views = append(views, PostgresViews()...)

	return views
}

func bucketView(instrumentName string, boundaries []float64) sdkmetric.View {
	return sdkmetric.NewView(
		sdkmetric.Instrument{Name: instrumentName},
		sdkmetric.Stream{
			Aggregation: sdkmetric.AggregationExplicitBucketHistogram{Boundaries: boundaries},
		},
	)
}
