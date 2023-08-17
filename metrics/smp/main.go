package main

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	ctx := context.Background()
	exp, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp)))
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(meterProvider)
	meter := otel.Meter("go.opentelemetry.io/otel/metric#MultiAsyncExample")

	// This is just a sample of memory stats to record from the Memstats
	cpuUsage, _ := meter.Int64ObservableGauge(
		"cpuUsage",
		api.WithDescription("CPU Usage in %"),
	)

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			memStats := &runtime.MemStats{}
			// This call does work
			runtime.ReadMemStats(memStats)
			o.ObserveInt64(cpuUsage,
				int64(60),
				api.WithAttributes(
					attribute.String("label", "value"),
					attribute.Bool("env-prod", true),
				),
			)
			return nil
		},
		cpuUsage,
	)
	if err != nil {
		fmt.Println("Failed to register callback")
		panic(err)
	}
}
