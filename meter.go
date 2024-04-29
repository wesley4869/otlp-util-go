package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// InitGlobalMeter initializes the global meter.
func InitGlobalMeter(opts ...InitOption) (*sdkmetric.MeterProvider, error) {
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

	provider, err := newMeterProvider(opt.ctx, opt.serviceName, opt.endpoint)
	if err != nil {
		return nil, fmt.Errorf("create meter provider. %w", err)
	}
	return provider, nil
}

// NewInt64Counter creates a new int64 counter.
func NewInt64Counter(name string, opts ...metric.Int64CounterOption) metric.Int64Counter {
	counter, _ := globalMeter.Int64Counter(name, opts...)
	return counter
}

// NewInt64UpDownCounter creates a new int64 up down counter.
func NewInt64UpDownCounter(name string, opts ...metric.Int64UpDownCounterOption) metric.Int64UpDownCounter {
	counter, _ := globalMeter.Int64UpDownCounter(name, opts...)
	return counter
}

// NewInt64Histogram creates a new int64 histogram.
func NewInt64Histogram(name string, opts ...metric.Int64HistogramOption) metric.Int64Histogram {
	histogram, _ := globalMeter.Int64Histogram(name, opts...)
	return histogram

}

// NewInt64ObservableCounter creates a new int64 observable counter.
func NewInt64ObservableCounter(name string, opts ...metric.Int64ObservableCounterOption) metric.Int64ObservableCounter {
	counter, _ := globalMeter.Int64ObservableCounter(name, opts...)
	return counter
}

// NewInt64ObservableUpDownCounter creates a new int64 observable up down counter.
func NewInt64ObservableUpDownCounter(name string, opts ...metric.Int64ObservableUpDownCounterOption) metric.Int64ObservableUpDownCounter {
	counter, _ := globalMeter.Int64ObservableUpDownCounter(name, opts...)
	return counter
}

// NewInt64ObservableGauge creates a new int64 observable gauge.
func NewInt64ObservableGauge(name string, opts ...metric.Int64ObservableGaugeOption) metric.Int64ObservableGauge {
	gauge, _ := globalMeter.Int64ObservableGauge(name, opts...)
	return gauge
}

// NewFloat64Counter creates a new float64 counter.
func NewFloat64Counter(name string, opts ...metric.Float64CounterOption) metric.Float64Counter {
	counter, _ := globalMeter.Float64Counter(name, opts...)
	return counter
}

// NewFloat64UpDownCounter creates a new float64 up down counter.
func NewFloat64UpDownCounter(name string, opts ...metric.Float64UpDownCounterOption) metric.Float64UpDownCounter {
	counter, _ := globalMeter.Float64UpDownCounter(name, opts...)
	return counter
}

// NewFloat64Histogram creates a new float64 histogram.
func NewFloat64Histogram(name string, opts ...metric.Float64HistogramOption) metric.Float64Histogram {
	histogram, _ := globalMeter.Float64Histogram(name, opts...)
	return histogram
}

// NewFloat64ObservableCounter creates a new float64 observable counter.
func NewFloat64ObservableCounter(name string, opts ...metric.Float64ObservableCounterOption) metric.Float64ObservableCounter {
	counter, _ := globalMeter.Float64ObservableCounter(name, opts...)
	return counter
}

// NewFloat64ObservableUpDownCounter creates a new float64 observable up down counter.
func NewFloat64ObservableUpDownCounter(name string, opts ...metric.Float64ObservableUpDownCounterOption) metric.Float64ObservableUpDownCounter {
	counter, _ := globalMeter.Float64ObservableUpDownCounter(name, opts...)
	return counter
}

// NewFloat64ObservableGauge creates a new float64 observable gauge.
func NewFloat64ObservableGauge(name string, opts ...metric.Float64ObservableGaugeOption) metric.Float64ObservableGauge {
	gauge, _ := globalMeter.Float64ObservableGauge(name, opts...)
	return gauge
}
