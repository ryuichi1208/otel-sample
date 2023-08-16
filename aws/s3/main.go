package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var tp *sdktrace.TracerProvider

func getList(ctx context.Context) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedConfigProfile("oneaws-colorme"),
		config.WithRegion("ap-northeast-1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	client := s3.NewFromConfig(cfg)

	for i := 0; i < 10; i++ {
		output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String("colorme-drone-cache"),
		})

		if err != nil {
			log.Fatal(err, output)
		}
		log.Println("All Objects in infytech2 bucket")
	}
}

const (
	service     = "trace-demo"
	environment = "production"
	id          = 1
)

func tracerProvider(url string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil
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
	otel.SetTracerProvider(tp)
	return nil
}

func initTracer() {
	var err error
	exp, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		fmt.Printf("failed to initialize stdout exporter %v\n", err)
		return
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)
}

func run() {
	tracer := tp.Tracer("example/aws/main")

	ctx := context.Background()
	defer func() { _ = tp.Shutdown(ctx) }()

	var span trace.Span
	ctx, span = tracer.Start(ctx, "AWS Example")
	defer span.End()
	getList(ctx)
}

func main() {
	err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	run()
}
