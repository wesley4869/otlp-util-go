package otlp_util

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var global_tracer trace.Tracer

// init() initialize global_tracer to a noop tracer.
func init() {
	global_tracer = otel.Tracer("noop")
}
