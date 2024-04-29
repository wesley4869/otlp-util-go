package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	otlp "github.com/bluexlab/otlp-util-go"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	ctx := context.Background()
	provider, err := otlp.InitGlobalMeter(
		otlp.WithContext(ctx),
		otlp.WithEndPoint("localhost:4317"),
		otlp.WithServiceName("otlp_util_example"),
		otlp.WithInSecure(),
		otlp.WithErrorHandler(func(err error) {
			log.Printf("OTLP error: %v", err)
		}),
	)
	if err != nil {
		log.Fatalf("failed to initialize OTLP exporter: %v", err)
	}
	defer func() { _ = provider.Shutdown(ctx) }()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	count := otlp.NewInt64Counter("meter.demo.count", metric.WithDescription("A count of things"))
	for {
		select {
		case <-t.C:
			count.Add(ctx, 1)
		case <-ctx.Done():
			return
		}
	}
}
