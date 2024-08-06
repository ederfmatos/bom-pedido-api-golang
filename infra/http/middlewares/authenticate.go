package middlewares

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/telemetry"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

func AuthenticateMiddleware(factory *factory.ApplicationFactory) echo.MiddlewareFunc {
	customerTokenManager := factory.TokenFactory.CustomerTokenManager
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, span := telemetry.StartSpan(c.Request().Context(), "AuthenticateMiddleware")
			c.SetRequest(c.Request().WithContext(ctx))
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				span.SetAttributes(
					attribute.String("customerId", "019078bc-cab8-789a-a1e7-4ba2a09561a6"),
					attribute.String("adminId", "019078bc-cab8-789a-a1e7-4ba2a09561a6"),
				)
				c.Set("customerId", "019078bc-cab8-789a-a1e7-4ba2a09561a6")
				c.Set("adminId", "019078bc-cab8-789a-a1e7-4ba2a09561a6")
				span.End()
				return next(c)
			}
			customerId, err := customerTokenManager.Decrypt(ctx, token)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
				span.End()
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			c.Set("customerId", customerId)
			c.Set("adminId", customerId)
			span.SetAttributes(
				attribute.String("customerId", "019078bc-cab8-789a-a1e7-4ba2a09561a6"),
				attribute.String("adminId", "019078bc-cab8-789a-a1e7-4ba2a09561a6"),
			)
			span.End()
			return next(c)
		}
	}
}
