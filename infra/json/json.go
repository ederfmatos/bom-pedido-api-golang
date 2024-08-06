package json

import (
	"bom-pedido-api/infra/telemetry"
	"context"
	"encoding/json"
	"io"
)

func Unmarshal(ctx context.Context, data []byte, v any) error {
	_, span := telemetry.StartSpan(ctx, "JSON.Unmarshal")
	defer span.End()
	return json.Unmarshal(data, v)
}

func Marshal(ctx context.Context, v any) ([]byte, error) {
	_, span := telemetry.StartSpan(ctx, "JSON.Marshal")
	defer span.End()
	return json.Marshal(v)
}

func Decode(ctx context.Context, r io.Reader, v any) error {
	_, span := telemetry.StartSpan(ctx, "JSON.Decode")
	defer span.End()
	return json.NewDecoder(r).Decode(v)
}

func Encode(ctx context.Context, w io.Writer, v any) error {
	_, span := telemetry.StartSpan(ctx, "JSON.Encode")
	defer span.End()
	return json.NewEncoder(w).Encode(v)
}
