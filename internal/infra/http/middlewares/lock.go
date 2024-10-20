package middlewares

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/telemetry"
	"context"
	"github.com/labstack/echo/v4"
	"time"
)

func LockByParam(name string, factory *factory.ApplicationFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			param := c.Param(name)
			ctx, span := telemetry.StartSpan(
				c.Request().Context(),
				"LockByParam",
				"param", name,
				"value", param,
			)
			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))
			lockKey, err := factory.Locker.Lock(c.Request().Context(), time.Minute, param)
			if err != nil {
				return err
			}
			go releaseOnContextDone(c.Request().Context(), factory, lockKey)
			defer factory.Locker.Release(context.Background(), lockKey)
			return next(c)
		}
	}
}

func LockByCustomerId(factory *factory.ApplicationFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Get(CustomerIdParam).(string)
			ctx, span := telemetry.StartSpan(c.Request().Context(), "LockByCustomerId", "customerId", id)
			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))
			lockKey, err := factory.Locker.Lock(c.Request().Context(), time.Minute, id)
			if err != nil {
				return err
			}
			go releaseOnContextDone(c.Request().Context(), factory, lockKey)
			defer factory.Locker.Release(context.Background(), lockKey)
			return next(c)
		}
	}
}

func releaseOnContextDone(ctx context.Context, factory *factory.ApplicationFactory, id string) {
	<-ctx.Done()
	factory.Locker.Release(context.Background(), id)
}
