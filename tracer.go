package otlp_util

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type InitOption func(*config)

type config struct {
	serviceName  string
	endPoint     string
	grpcInSecure bool
}

func InitGlobalTracer(opts ...InitOption) {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceNameKey.String(cfg.serviceName)),
	)
	if err != nil {
		panic(err)
	}

	options := make([]otlptracegrpc.Option, 0, 8)
	options = append(options, otlptracegrpc.WithEndpoint(cfg.endPoint))
	if cfg.grpcInSecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	client := otlptracegrpc.NewClient(options...)

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(err)
	}

	sp := sdktrace.NewBatchSpanProcessor(
		exporter,
	)

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithSpanProcessor(sp),
	)

	otel.SetTracerProvider(provider)
	tracer := otel.GetTracerProvider().Tracer(cfg.serviceName)
	global_tracer = tracer
}

func WithServiceName(name string) InitOption {
	return func(c *config) {
		c.serviceName = name
	}
}

func WithEndPoint(endPoint string) InitOption {
	return func(c *config) {
		c.endPoint = endPoint
	}
}

func WithInSecure() InitOption {
	return func(c *config) {
		c.grpcInSecure = true
	}
}

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return global_tracer.Start(ctx, spanName, opts...)
}
