package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
)

type (
	Span interface {
		End()
		RecordError(err error)
		SetAttribute(key, value string)
	}

	customSpan struct {
		span trace.Span
	}

	noneSpan struct {
	}
)

func (n noneSpan) End() {
}

func (n noneSpan) RecordError(_ error) {
}

func (n noneSpan) SetAttribute(_, _ string) {
}

func (s customSpan) SetAttribute(key, value string) {
	s.span.SetAttributes(attribute.String(key, value))
}

var tracer = otel.Tracer("bom-pedido-api")

func (s customSpan) RecordError(err error) {
	s.span.RecordError(err)
	s.span.SetStatus(1, err.Error())
}

func (s customSpan) End() {
	s.span.End()
}

var StartSpan func(ctx context.Context, spanName string, attributes ...string) (context.Context, Span)

func init() {
	if os.Getenv("ENV") == "production" {
		StartSpan = productionStartSpan
	} else {
		StartSpan = localStartSpan
	}
}

func productionStartSpan(ctx context.Context, spanName string, attributes ...string) (context.Context, Span) {
	attrs := make([]attribute.KeyValue, len(attributes)/2)
	index := 0
	for i := 0; i < len(attributes); i += 2 {
		attrs[index] = attribute.String(attributes[i], attributes[i+1])
		index++
	}
	ctx, span := tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
	return ctx, customSpan{span: span}
}

func localStartSpan(ctx context.Context, _ string, _ ...string) (context.Context, Span) {
	return ctx, noneSpan{}
}
