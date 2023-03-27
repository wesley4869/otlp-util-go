package otlp_util

import (
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var global_tracer trace.Tracer
var global_tracer_provider *sdktrace.TracerProvider

// init() initialize global_tracer to a noop tracer.
func init() {
	global_tracer = otel.Tracer("noop")
}
