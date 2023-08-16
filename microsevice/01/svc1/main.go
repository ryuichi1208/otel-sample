package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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
			semconv.ServiceNameKey.String("microservice-1"),
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
	resty := resty.New()

	router := gin.Default()
	router.Use(otelgin.Middleware("microservice-1"))
	{

		router.GET("/ping", func(c *gin.Context) {
			result := Result{}
			req := resty.R().SetHeader("Content-Type", "application/json")
			ctx := req.Context()
			span := trace.SpanFromContext(ctx)

			defer span.End()

			otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
			resp, _ := req.Get("http://localhost:8088/pong")

			json.Unmarshal([]byte(resp.String()), &result)
			c.IndentedJSON(200, gin.H{
				"message": result.Message,
			})
		})
	}
	router.Run(":8085")

}
