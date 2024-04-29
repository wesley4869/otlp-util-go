package main

import (
	"context"
	"fmt"
	"os"
	"time"

	otlp "github.com/bluexlab/otlp-util-go"
)

func level1(ctx context.Context) {
	ctx, span := otlp.Start(ctx, "level1")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
	level2(ctx)
	time.Sleep(100 * time.Millisecond)
}

func level2(ctx context.Context) {
	_, span := otlp.Start(ctx, "level2")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
}

func main() {
	provider, err := otlp.InitGlobalTracer(
		otlp.WithEndPoint("localhost:4317"),
		otlp.WithServiceName("otlp_util_example"),
		otlp.WithInSecure(),
		otlp.WithErrorHandler(func(err error) {
			fmt.Println(err.Error())
		}),
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { _ = provider.Shutdown(context.Background()) }()

	ctx := context.Background()

	for {
		level1(ctx)
		time.Sleep(1 * time.Second)
	}
}
