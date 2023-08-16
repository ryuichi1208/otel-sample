package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Result struct {
	Message string
}

func TracerProvider() (*tracesdk.TracerProvider, error) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("microservice-2"),
			semconv.ServiceVersionKey.String("0.0.1"),
			attribute.String("environment", "test"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	_, err := TracerProvider()
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.Use(otelgin.Middleware("microservice-2"))
	{
		router.GET("/pong", func(c *gin.Context) {
			ctx := c.Request.Context()
			span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)))
			defer span.End()

			c.IndentedJSON(200, gin.H{
				"message": "pong",
			})
		})
	}
	router.Run(":8088")

}
