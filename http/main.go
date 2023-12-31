package main

import (
	"context"
	"fmt"
	"log"
	"main/lib/http"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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
	// Create a new resource with service name and version
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(svcName),
		semconv.ServiceVersion(svcVersion),
		attribute.String("go", runtime.Version()),
	)

	// Create a new OTLP exporter over gRPC with no authentication and
	client := otlptracehttp.NewClient(otlptracehttp.WithInsecure())
	// Create a new trace exporter
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	// Create a new trace provider with the exporter
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(
			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio)),
		),
		sdktrace.WithResource(resource),
	)

	// Register the trace provider with the global tracer provider
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

func main() {
	ctx := context.Background()
	shutdown, err := installExportPipeline(ctx, "otlptrace-example", "0.0.1", 1.0)
	if err != nil {
		log.Fatal(err)
	}
	ctx, span := tracer.Start(ctx, "main")
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	c := http.New(ctx, tracer)
	c.Do(ctx)
	span.End()
}
