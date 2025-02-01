package lock

import (
	"bom-pedido-api/internal/domain/errors"
	"context"
)

var ResourceLockedError = errors.New("Resource locked")

type Locker interface {
	LockFunc(ctx context.Context, key string, lockedFunc func()) error
	Lock(ctx context.Context, key ...string) (string, error)
	Release(ctx context.Context, key string)
}
