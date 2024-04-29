package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
)

type option struct {
	ctx          context.Context
	serviceName  string
	endpoint     string
	errorHandler otel.ErrorHandler
	grpcInsecure bool
}

type InitOption func(*option)

// WithContext sets the context
func WithContext(ctx context.Context) InitOption {
	return func(o *option) {
		o.ctx = ctx
	}
}

// WithServiceName sets the service name.
func WithServiceName(serviceName string) InitOption {
	return func(o *option) {
		o.serviceName = serviceName
	}
}

// WithEndPoint sets the endpoint.
func WithEndPoint(endpoint string) InitOption {
	return func(o *option) {
		o.endpoint = endpoint
	}
}

// WithErrorHandler sets the error handler.
func WithErrorHandler(handler otel.ErrorHandlerFunc) InitOption {
	return func(o *option) {
		o.errorHandler = handler
	}
}

// WithInSecure sets the gRPC inscure mode
func WithInSecure() InitOption {
	return func(o *option) {
		o.grpcInsecure = true
	}
}

func (o *option) Validate() error {
	if o.endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	if o.serviceName == "" {
		return fmt.Errorf("service name is required")
	}
	return nil
}
