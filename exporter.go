package otlp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"google.golang.org/grpc"
)

type ShutdownHandler interface {
	Shutdown(ctx context.Context) error
}

type ShutdownFunc func(ctx context.Context) error

func (f ShutdownFunc) Shutdown(ctx context.Context) error {
	return f(ctx)
}

// InitExporter initializes the global tracer and meter.
func InitExporter(opts ...InitOption) (ShutdownHandler, error) {
	opt := option{ctx: context.Background()}
	for _, o := range opts {
		o(&opt)
	}
	if opt.Validate() != nil {
		return nil, fmt.Errorf("validate option. %w", opt.Validate())
	}

	// Set error handler
	if opt.errorHandler != nil {
		otel.SetErrorHandler(opt.errorHandler)
	}

	// Create meter exporter
	metricProvider, err := newMeterProvider(opt.ctx, opt.serviceName, opt.endpoint)
	if err != nil {
		return nil, fmt.Errorf("create meter provider. %w", err)
	}

	// Create trace exporter
	traceProvider, err := newTraceProvider(opt.ctx, opt.serviceName, opt.endpoint)
	if err != nil {
		return nil, fmt.Errorf("create trace provider. %w", err)
	}

	shutdown := ShutdownFunc(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = metricProvider.ForceFlush(ctx)
		if err := metricProvider.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
			otel.Handle(err)
		}

		_ = traceProvider.ForceFlush(ctx)
		if err := traceProvider.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
			otel.Handle(err)
		}

		return nil
	})
	return shutdown, nil
}

func getResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithHost(),
		resource.WithHostID(),
		resource.WithContainer(),
		resource.WithContainerID(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
}

func newMeterProvider(ctx context.Context, serviceName string, endpoint string) (*sdkmetric.MeterProvider, error) {
	res, err := getResource(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("create resource. %w", err)
	}

	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("create metric exporter. %w", err)
	}

	pr := sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(5*time.Second))
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(pr),
	)
	otel.SetMeterProvider(provider)

	// Set global meter
	globalMeter = provider.Meter(serviceName)

	return provider, nil
}

func newTraceProvider(ctx context.Context, serviceName string, endpoint string) (*sdktrace.TracerProvider, error) {
	res, err := getResource(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("create resource. %w", err)
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	traceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("create trace exporter. %w", err)
	}

	sp := sdktrace.NewBatchSpanProcessor(traceExporter)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sp),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)

	// Set global tracer
	globalTracer = provider.Tracer(serviceName)

	return provider, nil
}
