package http

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	t trace.Tracer
}

func New(ctx context.Context, t trace.Tracer) Client {
	// installExportPipeline(ctx)
	return Client{
		t: t,
	}
}

func (c Client) add(ctx context.Context, x, y int64) int64 {
	var span trace.Span
	_, span = c.t.Start(ctx, "Addition2")
	defer span.End()

	return x + y
}

func (c Client) Do(ctx context.Context) int64 {
	var span trace.Span
	_, span = c.t.Start(ctx, "Addition1")
	defer span.End()

	c.add(ctx, 1, 2)

	fmt.Println("Hello, World!")
	return 10
}
