// Package tracing provides OpenTelemetry tracing setup.
package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer initializes and returns an OTLP tracer provider.
func InitTracer() *sdktrace.TracerProvider {
	// OTLP HTTP exporter
	exp, err := otlptracehttp.New(context.Background())
	if err != nil {
		log.Fatal("Failed to create OTLP exporter:", err)
	}

	// Tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	return tp
}
