package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service     = "trace-demo"
	environment = "production"
	id          = 1
)

func tracerProvider(url string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-main")

	ctx, span := tr.Start(ctx, "foo")
	defer span.End()

	bar(ctx)

	httpRequest(ctx)
}

func bar(ctx context.Context) {
	// Use the global TracerProvider.
	tr := otel.Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	// Do bar...
}

func sleep2(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "sleep")
	defer span.End()
	time.Sleep(2 * time.Second)
}

func sleep(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "sleep")
	defer span.End()
	time.Sleep(3 * time.Second)
	sleep2(ctx, tracer)
}

func httpRequest(ctx context.Context) error {
	tracer := otel.Tracer("component-bar2")
	var span trace.Span
	ctx, span = tracer.Start(ctx, "httpRequest")
	defer span.End()

	resp, err := otelhttp.Get(ctx, "https://google.com/")
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	httpRequest2(ctx)
	return nil
}

func httpRequest2(ctx context.Context) error {
	tracer := otel.Tracer("component-bar3")
	var span trace.Span
	ctx, span = tracer.Start(ctx, "httpRequest2")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "example-service/1.0.0")
	cli := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
