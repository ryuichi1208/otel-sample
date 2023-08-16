package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()
	shutdown, err := installExportPipeline(ctx, "otlptrace-example", "0.0.1", 1.0)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	h := newHandler()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.hello))
	http.ListenAndServe(":8001", otelhttp.NewHandler(mux, "api"))
}

type handler struct {
	cli http.Client
}

func newHandler() *handler {
	htp := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	hc := http.Client{
		Transport: otelhttp.NewTransport(htp), // instrument http.Transport
		Timeout:   60 * time.Second,
	}
	return &handler{
		cli: hc,
	}
}

func (h *handler) hello(w http.ResponseWriter, r *http.Request) {
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	ctx := r.Context()
	_, span := tracer.Start(ctx, "op2")
	defer span.End()
	clientTrace := otelhttptrace.NewClientTrace(ctx)
	ctx = httptrace.WithClientTrace(ctx, clientTrace)
	hreq, err := http.NewRequestWithContext(ctx, "GET", "http://httpbin.org/delay/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := h.cli.Do(hreq)
	if err != nil {
		log.Fatal(err)
	}
	// span ends after resp.Body.Close.
	resp.Body.Close()
}
