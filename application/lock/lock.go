package lock

import (
	"context"
	"time"
)

type Locker interface {
	LockFunc(ctx context.Context, key string, ttl time.Duration, lockedFunc func()) error
	Lock(ctx context.Context, ttl time.Duration, key ...string) (string, error)
	Release(ctx context.Context, key string) error
}
