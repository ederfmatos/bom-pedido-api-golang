package telemetry

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"os"
)

type (
	TracerProvider interface {
		Close() error
	}

	OtelTracerProvider struct {
		provider *trace.TracerProvider
	}
)

func NewTracerProvider() (TracerProvider, error) {
	var (
		ctx         = context.Background()
		serviceName = "bom-pedido-api"
		exporterURL = os.Getenv("OTEL_ENDPOINT_EXPORTER")
	)

	res, err := resource.New(ctx, resource.WithAttributes(attribute.Key("service.name").String(serviceName)))
	if err != nil {
		return nil, fmt.Errorf("create resource: %v", err)
	}

	headers := map[string]string{
		"content-type": "application/json",
	}

	tracerClient := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(exporterURL),
		otlptracehttp.WithHeaders(headers),
		otlptracehttp.WithInsecure(),
	)
	traceExporter, err := otlptrace.New(ctx, tracerClient)
	if err != nil {
		return nil, fmt.Errorf("create tracer exporter: %v", err)
	}

	spanProcessor := trace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(spanProcessor),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return &OtelTracerProvider{provider: tracerProvider}, nil
}

func (p OtelTracerProvider) Close() error {
	return p.provider.Shutdown(context.Background())
}
