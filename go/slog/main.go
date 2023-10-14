package main

import (
	"context"
	"log/slog"
	"os"
)

type MyHandler struct {
	slog.Handler
}

type traceIDKey struct {
	ID string
}

var tidKey = traceIDKey{
	ID: "traceID",
}

func (h MyHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("traceID", ctx.Value(tidKey).(string)))
	return h.Handler.Handle(ctx, r)
}

func main() {
	myHandler := MyHandler{
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	}
	logger := slog.New(myHandler)

	ctx := context.WithValue(context.Background(), tidKey, "12345")
	logger.InfoContext(ctx, "Hello World1", "foo1", "bar1")

	slog.SetDefault(logger)
	slog.InfoContext(ctx, "Hello World2", "foo2", "bar2")
}
