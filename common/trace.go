package common

import (
	"context"
	"log"
	"os"

	"github.com/pandada8/otel-config-go/otelconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	ServiceName = os.Getenv("OTEL_SERVICE_NAME")
	Meter       = otel.Meter(ServiceName)
)

func InitTracer() func() {
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry()
	if err != nil {
		log.Fatalf("error setup otel sdk: %s", err)
	}
	return otelShutdown
}

func GetTraceId(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	return span.SpanContext().TraceID().String()
}
