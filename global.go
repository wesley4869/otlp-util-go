package otlp

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	globalMeter  = otel.Meter("noop")
	globalTracer = otel.Tracer("noop")
)

// SetGlobalMeter sets the global meter.
func SetGlobalMeter(meter metric.Meter) {
	globalMeter = meter
}

// GlobalMeter returns the global meter.
func GlobalMeter() metric.Meter {
	return globalMeter
}

// SetGlobalTracer sets the global tracer.
func SetGlobalTracer(tracer trace.Tracer) {
	globalTracer = tracer
}

// GlobalTracer returns the global tracer.
func GlobalTracer() trace.Tracer {
	return globalTracer
}
