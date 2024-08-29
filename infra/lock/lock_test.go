package lock

import (
	"bom-pedido-api/infra/test"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisLocker(t *testing.T) {
	container := test.NewContainer()
	defer container.Down()
	redisClient := container.RedisClient

	locker := NewRedisLocker(redisClient)

	t.Run("Lock and Release", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key"
		ttl := 10 * time.Second

		_, err := locker.Lock(ctx, ttl, key)
		require.NoError(t, err, "failed to lock:", err)

		err = locker.Release(ctx, key)
		if err != nil {
			t.Fatalf("failed to release lock: %s", err)
		}
	})

	t.Run("LockFunc", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_func"
		ttl := 10 * time.Second

		called := false
		err := locker.LockFunc(ctx, key, ttl, func() {
			called = true
		})
		require.NoError(t, err, "failed to lock:", err)

		if !called {
			t.Fatal("locked function was not called")
		}
	})

	t.Run("Lock when already locked", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_locked"
		ttl := 10 * time.Second

		_, err := locker.Lock(ctx, ttl, key)
		require.NoError(t, err, "failed to lock:", err)

		_, err = locker.Lock(ctx, ttl, key)
		if err == nil {
			t.Fatal("expected lock to fail but it succeeded")
		}

		err = locker.Release(ctx, key)
		if err != nil {
			t.Fatalf("failed to release lock: %s", err)
		}
	})

	t.Run("Lock with expired TTL", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_expired"
		ttl := 2 * time.Second

		_, err := locker.Lock(ctx, ttl, key)
		require.NoError(t, err, "failed to lock:", err)

		time.Sleep(3 * time.Second)

		_, err = locker.Lock(ctx, ttl, key)
		if err != nil {
			t.Fatalf("failed to re-lock after TTL expired: %s", err)
		}

		err = locker.Release(ctx, key)
		if err != nil {
			t.Fatalf("failed to release lock: %s", err)
		}
	})

	t.Run("Lock with canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		key := "test_key_cancel"
		ttl := 10 * time.Second

		_, err := locker.Lock(ctx, ttl, key)
		require.NoError(t, err, "failed to lock:", err)

		cancel()

		time.Sleep(1 * time.Second)

		locked, err := redisClient.Get(context.Background(), key).Result()
		if !errors.Is(err, redis.Nil) && err != nil {
			t.Fatalf("unexpected error getting key: %s", err)
		}

		if locked != "" {
			t.Fatal("expected lock to be released after context cancellation")
		}
	})
}
