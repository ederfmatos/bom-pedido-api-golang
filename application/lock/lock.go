package lock

import (
	"bom-pedido-api/domain/errors"
	"context"
	"time"
)

var ResourceLockedError = errors.New("Resource locked")

type Locker interface {
	LockFunc(ctx context.Context, key string, ttl time.Duration, lockedFunc func()) error
	Lock(ctx context.Context, ttl time.Duration, key ...string) (string, error)
	Release(ctx context.Context, key string)
}
