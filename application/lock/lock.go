package lock

import (
	"context"
	"time"
)

type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration, lockedFunc func()) error
}
