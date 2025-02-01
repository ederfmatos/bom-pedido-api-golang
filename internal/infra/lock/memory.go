package lock

import (
	"bom-pedido-api/internal/application/lock"
	"context"
	"strings"
	"sync"
	"time"
)

type (
	memoryLocker struct {
		mutex sync.Mutex
		locks map[string]time.Time
	}
)

func NewMemoryLocker() lock.Locker {
	return &memoryLocker{
		locks: make(map[string]time.Time),
	}
}

func (l *memoryLocker) Lock(_ context.Context, key ...string) (string, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	lockKey := strings.Join(key, "")

	expiration, exists := l.locks[lockKey]
	if exists && time.Now().Before(expiration) {
		return "", lock.ResourceLockedError
	}

	return lockKey, nil
}

func (l *memoryLocker) LockFunc(ctx context.Context, key string, lockedFunc func()) error {
	lockKey, err := l.Lock(ctx, key)
	if err != nil {
		return err
	}
	lockedFunc()
	l.Release(context.Background(), lockKey)
	return nil
}

func (l *memoryLocker) Release(_ context.Context, key string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, exists := l.locks[key]
	if !exists {
		return
	}

	delete(l.locks, key)
}
