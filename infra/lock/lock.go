package lock

import (
	"bom-pedido-api/application/lock"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
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

func (l *redisLocker) Lock(ctx context.Context, key string, ttl time.Duration, lockedFunc func()) error {
	locked, err := l.client.SetNX(ctx, key, "", ttl).Result()
	if err != nil {
		return err
	}
	if !locked {
		return fmt.Errorf("redis lock failed")
	}
	go func() {
		select {
		case <-ctx.Done():
			_ = l.unlock(ctx, &key)
		}
	}()
	lockedFunc()
	_ = l.unlock(ctx, &key)
	return nil
}

func (l *redisLocker) unlock(ctx context.Context, id *string) error {
	releaseScript := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	_, err := l.client.Eval(ctx, releaseScript, []string{*id}, "").Bool()
	return err
}
