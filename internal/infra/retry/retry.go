package retry

import (
	"bom-pedido-api/pkg/log"
	"bom-pedido-api/pkg/telemetry"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"
)

func Func(ctx context.Context, maxRetries int, initialBackoff, maxBackoff time.Duration, operation func(context.Context) error) error {
	err := operation(ctx)
	if err == nil {
		return nil
	}
	for attempt := 2; attempt <= maxRetries; attempt++ {
		err = telemetry.StartSpanReturningError(ctx, fmt.Sprintf("Retry::%v", attempt), func(ctx context.Context) error {
			err = operation(ctx)
			if err == nil {
				return nil
			}
			backoff := time.Duration(math.Min(float64(maxBackoff), float64(initialBackoff)*math.Pow(2, float64(attempt-1))))
			log.Warn("Attempt failed. Retrying...", "attempt", attempt, "backoff", backoff.String(), "err", err)
			time.Sleep(backoff)
			return err
		}, "attempt", strconv.Itoa(attempt))

		if err == nil {
			return nil
		}
	}
	log.Error("operation failed", err, "maxRetries", maxRetries, err)
	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, err)
}
