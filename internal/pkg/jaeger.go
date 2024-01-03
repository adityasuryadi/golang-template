package pkg

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("order-service")

type JaegerTracer struct {
	Exporter *otlptrace.Exporter
	Tracer   trace.Tracer
}

func NewJaegerTracer(exporter *otlptrace.Exporter) *JaegerTracer {
	return &JaegerTracer{
		Exporter: exporter,
		Tracer:   tracer,
	}
}

func (e *JaegerTracer) Trace(endpoint string) *sdktrace.TracerProvider {
	exporter := e.Exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(endpoint),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
