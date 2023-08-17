package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
)

var (
	tracer = otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
)

func installExportPipeline(ctx context.Context, svcName, svcVersion string, ratio float64) (func(context.Context) error, error) {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	// Create a new resource with service name and version
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(svcName),
		semconv.ServiceVersion(svcVersion),
	)

	// Create a new OTLP exporter over gRPC with no authentication and
	client := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
	)

	// Create a new trace exporter
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	// Create a new trace provider with the exporter
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(
			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.5)),
		),
		sdktrace.WithResource(resource),
	)

	// Register the trace provider with the global tracer provider
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}
