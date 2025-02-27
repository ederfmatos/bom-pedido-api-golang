package lock

import (
	"bom-pedido-api/internal/application/lock"
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
)

type (
	redisLocker struct {
		client *redis.Client
	}
)

func NewRedisLocker(client *redis.Client) lock.Locker {
	return &redisLocker{client: client}
}

func (l *redisLocker) Lock(ctx context.Context, key ...string) (string, error) {
	lockKey := strings.Join(key, "")
	locked, err := l.client.SetNX(ctx, lockKey, "", -1).Result()
	if err != nil {
		return "", err
	}
	if !locked {
		return "", lock.ResourceLockedError
	}
	go func() {
		<-ctx.Done()
		l.Release(context.Background(), lockKey)
	}()
	return lockKey, nil
}
func (l *redisLocker) LockFunc(ctx context.Context, key string, lockedFunc func()) error {
	lockKey, err := l.Lock(ctx, key)
	if err != nil {
		return err
	}
	defer l.Release(context.Background(), lockKey)
	lockedFunc()
	return nil
}

func (l *redisLocker) Release(ctx context.Context, key string) {
	releaseScript := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	_, _ = l.client.Eval(ctx, releaseScript, []string{key}, "").Bool()
}
