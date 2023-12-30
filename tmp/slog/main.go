package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

func newLogger(traceID string) *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			},
		),
	).With("trace_id", traceID)
	return logger
}

var logger *slog.Logger

func add(a, b int) int {
	logger.Info("add", slog.Int("a", a), slog.Int("b", b))
	return a + b
}

func main() {
	logger = newLogger("123456")
	logger.Info("Hello World!!!")
	add(1, 2)
	main2()
}

func randomSleep(i int) {
	n := rand.Intn(3) + 1
	fmt.Printf("%d: Sleeping %d seconds...\n", i, n)
	time.Sleep(time.Duration(n) * time.Second)
}

func main2() {
	eg, ctx := errgroup.WithContext(context.Background())

	for i := 0; i < 4; i++ {
		i := i
		eg.Go(func() error {
			fmt.Println("Worker Index:", i)

			randomSleep(i)

			select {
			case <-ctx.Done():
				fmt.Println("Canceled:", i)
				return nil
			default:
				if i > 1 {
					fmt.Println("Fire Error:", i)
					return fmt.Errorf("Error: %d", i)
				}
				fmt.Println("Done Worker:", i)
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal("Catched err: ", err)
	}

}
