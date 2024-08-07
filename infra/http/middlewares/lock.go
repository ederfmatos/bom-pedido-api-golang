package middlewares

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/telemetry"
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
			err := factory.Locker.Lock(c.Request().Context(), param, time.Minute)
			if err != nil {
				return err
			}
			go ReleaseOnContextDone(c.Request().Context(), factory, param)
			defer factory.Locker.Release(context.Background(), param)
			return next(c)
		}
	}
}

func LockByCustomerId(factory *factory.ApplicationFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Get("customerId").(string)
			ctx, span := telemetry.StartSpan(c.Request().Context(), "LockByCustomerId", "customerId", id)
			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))
			err := factory.Locker.Lock(c.Request().Context(), id, time.Minute)
			if err != nil {
				return err
			}
			go ReleaseOnContextDone(c.Request().Context(), factory, id)
			defer factory.Locker.Release(context.Background(), id)
			return next(c)
		}
	}
}

func ReleaseOnContextDone(ctx context.Context, factory *factory.ApplicationFactory, id string) {
	select {
	case <-ctx.Done():
		factory.Locker.Release(context.Background(), id)
	}
}
