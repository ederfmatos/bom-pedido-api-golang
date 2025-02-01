package retry

import (
	"bom-pedido-api/internal/infra/telemetry"
	"bom-pedido-api/pkg/log"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"
)

func Func(ctx context.Context, maxRetries int, initialBackoff, maxBackoff time.Duration, operation func(context.Context) error) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, span := telemetry.StartSpan(ctx, fmt.Sprintf("Retry::%v", attempt), "attempt", strconv.Itoa(attempt))
		err = operation(ctx)
		if err == nil {
			span.End()
			return nil
		}
		span.RecordError(err)
		backoff := time.Duration(math.Min(float64(maxBackoff), float64(initialBackoff)*math.Pow(2, float64(attempt-1))))
		log.Warn("Attempt failed. Retrying...", "attempt", attempt, "backoff", backoff)
		time.Sleep(backoff)
		span.End()
	}
	log.Error("operation failed", err, "maxRetries", maxRetries, err)
	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, err)
}
