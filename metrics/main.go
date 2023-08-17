package main

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	ctx := context.Background()
	exp, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithInsecure())
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
}

func f() {
	meter := otel.Meter("go.opentelemetry.io/otel/metric#MultiAsyncExample")

	// This is just a sample of memory stats to record from the Memstats
	heapAlloc, _ := meter.Int64ObservableUpDownCounter("heapAllocs")
	gcCount, _ := meter.Int64ObservableCounter("gcCount")

	_, err := meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			memStats := &runtime.MemStats{}
			// This call does work
			runtime.ReadMemStats(memStats)

			o.ObserveInt64(heapAlloc, int64(memStats.HeapAlloc))
			o.ObserveInt64(gcCount, int64(memStats.NumGC))

			return nil
		},
		heapAlloc,
		gcCount,
	)
	if err != nil {
		fmt.Println("Failed to register callback")
		panic(err)
	}
}
