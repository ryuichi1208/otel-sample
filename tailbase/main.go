package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var tracer = otel.Tracer("demo-server")

func initProvider(ctx context.Context) func() {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("demo-server"),
		),
	)
	handleErr(err, "failed to create resource")

	otelAgentAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelAgentAddr = "0.0.0.0:4317"
	}

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	traceExp, err := otlptrace.New(ctx, traceClient)
	handleErr(err, "Failed to create the collector trace exporter")

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func main() {
	ctx := context.Background()

	shutdown := initProvider(ctx)
	defer shutdown()

	hello := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Print("/hello")
		time.Sleep(1 * time.Second)
		if _, err := w.Write([]byte("hello")); err != nil {
			http.Error(w, "write operation failed.", http.StatusInternalServerError)
		}
	})

	world := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Print("/world")
		if _, err := w.Write([]byte("world")); err != nil {
			http.Error(w, "write operation failed.", http.StatusInternalServerError)
		}
	})

	google := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Print("/api/google")
		if err := httpRequest(req.Context(), "https://www.google.com"); err != nil {
			log.Printf("error: %v", err)
		}
		if _, err := w.Write([]byte("success to request Google")); err != nil {
			http.Error(w, "write operation failed.", http.StatusInternalServerError)
		}
	})

	yahoo := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Print("/api/yahoo")
		if err := httpRequest(req.Context(), "https://www.yahoo.co.jp"); err != nil {
			log.Printf("error: %v", err)
		}
		if _, err := w.Write([]byte("success to request Yahoo")); err != nil {
			http.Error(w, "write operation failed.", http.StatusInternalServerError)
		}
	})

	mux := http.NewServeMux()
	mux.Handle("/hello", otelhttp.NewHandler(hello, "/hello"))
	mux.Handle("/world", otelhttp.NewHandler(world, "/world"))
	mux.Handle("/api/google", otelhttp.NewHandler(google, "/api/google"))
	mux.Handle("/api/yahoo", otelhttp.NewHandler(yahoo, "/api/yahoo"))

	fmt.Println("Server is running...")
	server := &http.Server{
		Addr: ":8080",
		Handler: otelhttp.NewHandler(mux, "server",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		),
		ReadHeaderTimeout: 20 * time.Second,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		handleErr(err, "server failed to serve")
	}
}

func httpRequest(ctx context.Context, url string) error {
	var span trace.Span
	ctx, span = tracer.Start(ctx, "httpRequest")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
