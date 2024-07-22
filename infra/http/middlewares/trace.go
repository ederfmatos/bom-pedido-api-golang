package middlewares

import (
	"bom-pedido-api/infra/telemetry"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TraceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := telemetry.StartSpan(
			c.Request().Context(),
			c.Request().Method+" "+c.Request().URL.Path,
			"http.method", c.Request().Method,
			"http.url", c.Request().URL.Path,
			"http.host", c.Request().URL.Host,
		)
		defer span.End()
		c.SetRequest(c.Request().WithContext(ctx))
		err := next(c)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return err
		}
		span.SetAttributes(attribute.Int("http.status_code", c.Response().Status))
		span.SetStatus(codes.Ok, "")
		return nil
	}
}
