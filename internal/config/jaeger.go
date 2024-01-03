package config

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func NewJaegerTracer(viper *viper.Viper, log *logrus.Logger) *otlptrace.Exporter {
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(viper.GetString("jaeger.host")))
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	return exporter
}
