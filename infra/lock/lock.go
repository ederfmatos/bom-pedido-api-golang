package lock

import (
	"bom-pedido-api/application/lock"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type (
	redisLocker struct {
		client *redis.Client
	}
)

func NewRedisLocker(client *redis.Client) lock.Locker {
	return &redisLocker{client: client}
}

func (l *redisLocker) Lock(ctx context.Context, ttl time.Duration, key ...string) (string, error) {
	lockKey := strings.Join(key, "")
	locked, err := l.client.SetNX(ctx, lockKey, "", ttl).Result()
	if err != nil {
		return "", err
	}
	if !locked {
		return "", fmt.Errorf("resource locked")
	}
	go func() {
		select {
		case <-ctx.Done():
			_ = l.Release(context.Background(), lockKey)
		}
	}()
	return lockKey, nil
}
func (l *redisLocker) LockFunc(ctx context.Context, key string, ttl time.Duration, lockedFunc func()) error {
	lockKey, err := l.Lock(ctx, ttl, key)
	if err != nil {
		return err
	}
	defer l.Release(context.Background(), lockKey)
	lockedFunc()
	return nil
}

func (l *redisLocker) Release(ctx context.Context, key string) error {
	releaseScript := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	_, err := l.client.Eval(ctx, releaseScript, []string{key}, "").Bool()
	return err
}
