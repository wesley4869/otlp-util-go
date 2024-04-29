package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// InitGlobalTracer initializes the global tracer.
func InitGlobalTracer(opts ...InitOption) (*sdktrace.TracerProvider, error) {
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

	provider, err := newTraceProvider(opt.ctx, opt.serviceName, opt.endpoint)
	if err != nil {
		return nil, fmt.Errorf("create trace provider. %w", err)
	}
	return provider, nil
}

// Start starts a new span.
func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return globalTracer.Start(ctx, spanName, opts...)
}
