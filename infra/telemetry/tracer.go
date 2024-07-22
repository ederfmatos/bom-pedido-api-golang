package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("bom-pedido-api")

func StartSpan(ctx context.Context, spanName string, attributes ...string) (context.Context, trace.Span) {
	attrs := make([]attribute.KeyValue, len(attributes)/2)
	index := 0
	for i := 0; i < len(attributes); i += 2 {
		attrs[index] = attribute.String(attributes[i], attributes[i+1])
		index++
	}
	return tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}
