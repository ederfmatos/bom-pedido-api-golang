package retry

import (
	"fmt"
	"log/slog"
	"math"
	"time"
)

func Func(maxRetries int, initialBackoff, maxBackoff time.Duration, operation func() error) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = operation()
		if err == nil {
			return nil
		}
		backoff := time.Duration(math.Min(float64(maxBackoff), float64(initialBackoff)*math.Pow(2, float64(attempt-1))))
		slog.Warn(fmt.Sprintf("Attempt %d failed; retrying in %s...", attempt, backoff))
		time.Sleep(backoff)
	}
	slog.Error(fmt.Sprintf("operation failed after %d attempts: %f", maxRetries, err))
	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, err)
}
