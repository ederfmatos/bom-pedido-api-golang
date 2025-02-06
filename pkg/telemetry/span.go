package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("bom-pedido-api")

func StartSpan[T any](ctx context.Context, spanName string, f func(ctx context.Context) (T, error), attributes ...string) (T, error) {
	ctx, span := startSpan(ctx, spanName, attributes...)
	output, err := f(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	span.End()
	return output, err
}

func StartSpanReturningError(ctx context.Context, spanName string, f func(ctx context.Context) error, attributes ...string) error {
	ctx, span := startSpan(ctx, spanName, attributes...)
	err := f(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(1, err.Error())
	}
	span.End()
	return err
}

func startSpan(ctx context.Context, spanName string, attributes ...string) (context.Context, trace.Span) {
	attrs := make([]attribute.KeyValue, len(attributes)/2)
	index := 0
	for i := 0; i < len(attributes); i += 2 {
		attrs[index] = attribute.String(attributes[i], attributes[i+1])
		index++
	}

	if tenantID := ctx.Value("TENANT_ID"); tenantID != nil {
		attrs = append(attrs, attribute.String("tenant.id", tenantID.(string)))
	}

	if userID := ctx.Value("USER_ID"); userID != nil {
		attrs = append(attrs, attribute.String("user.id", userID.(string)))
	}
	return tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}
