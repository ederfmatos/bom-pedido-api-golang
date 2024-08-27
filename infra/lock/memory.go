package lock

import (
	"bom-pedido-api/application/lock"
	"bom-pedido-api/domain/errors"
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

func (l *memoryLocker) Lock(_ context.Context, _ time.Duration, key ...string) (string, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	lockKey := strings.Join(key, "")

	expiration, exists := l.locks[lockKey]
	if exists && time.Now().Before(expiration) {
		return "", errors.New("lock already acquired")
	}

	return lockKey, nil
}

func (l *memoryLocker) LockFunc(ctx context.Context, key string, ttl time.Duration, lockedFunc func()) error {
	lockKey, err := l.Lock(ctx, ttl, key)
	if err != nil {
		return err
	}
	lockedFunc()
	return l.Release(context.Background(), lockKey)
}

func (l *memoryLocker) Release(_ context.Context, key string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, exists := l.locks[key]
	if !exists {
		return errors.New("lock not found")
	}

	delete(l.locks, key)
	return nil
}