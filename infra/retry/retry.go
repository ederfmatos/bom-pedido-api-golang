package retry

import (
	"fmt"
	"log/slog"
	"math"
	"time"
)

type Retry struct {
	maxRetries     int
	initialBackoff time.Duration
	maxBackoff     time.Duration
}

func NewRetry(maxRetries int, initialBackoff, maxBackoff time.Duration) *Retry {
	return &Retry{
		maxRetries:     maxRetries,
		initialBackoff: initialBackoff,
		maxBackoff:     maxBackoff,
	}
}

func (retry *Retry) Execute(operation func() error) error {
	var err error
	for attempt := 1; attempt <= retry.maxRetries; attempt++ {
		err = operation()
		if err == nil {
			return nil
		}
		backoff := time.Duration(math.Min(float64(retry.maxBackoff), float64(retry.initialBackoff)*math.Pow(2, float64(attempt-1))))
		slog.Warn(fmt.Sprintf("Attempt %d failed; retrying in %s...", attempt, backoff))
		time.Sleep(backoff)
	}
	slog.Error(fmt.Sprintf("operation failed after %d attempts: %f", retry.maxRetries, err))
	return fmt.Errorf("operation failed after %d attempts: %w", retry.maxRetries, err)
}
