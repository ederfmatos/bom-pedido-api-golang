package lock

import (
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
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

		_, err := locker.Lock(ctx, key)
		require.NoError(t, err, "failed to lock:", err)

		locker.Release(ctx, key)
	})

	t.Run("LockFunc", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_func"

		called := false
		err := locker.LockFunc(ctx, key, func() {
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

		_, err := locker.Lock(ctx, key)
		require.NoError(t, err, "failed to lock:", err)

		_, err = locker.Lock(ctx, key)
		if err == nil {
			t.Fatal("expected lock to fail but it succeeded")
		}

		locker.Release(ctx, key)
	})

	t.Run("Lock with expired TTL", func(t *testing.T) {
		t.Skip()

		ctx := context.Background()
		key := "test_key_expired"
		_, err := locker.Lock(ctx, key)
		require.NoError(t, err, "failed to lock:", err)

		time.Sleep(3 * time.Second)

		_, err = locker.Lock(ctx, key)
		if err != nil {
			t.Fatalf("failed to re-lock after TTL expired: %s", err)
		}

		locker.Release(ctx, key)
	})

	t.Run("Lock with canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		key := "test_key_cancel"

		_, err := locker.Lock(ctx, key)
		require.NoError(t, err, "failed to lock:", err)

		cancel()

		locked, err := redisClient.Get(context.Background(), key).Result()
		if !errors.Is(err, redis.Nil) && err != nil {
			t.Fatalf("unexpected error getting key: %s", err)
		}

		if locked != "" {
			t.Fatal("expected lock to be released after context cancellation")
		}
	})
}
