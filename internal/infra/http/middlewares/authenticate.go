package middlewares

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/infra/telemetry"
	"bom-pedido-api/internal/infra/tenant"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	AdminIdParam    = "adminId"
	CustomerIdParam = "customerId"
)

func AuthenticateMiddleware(factory *factory.ApplicationFactory) echo.MiddlewareFunc {
	customerTokenManager := factory.TokenFactory.TokenManager
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return next(c)
			}
			ctx, span := telemetry.StartSpan(c.Request().Context(), "AuthenticateMiddleware")
			tokenData, err := customerTokenManager.Decrypt(ctx, strings.ReplaceAll(token, "Bearer ", ""))
			if err != nil {
				span.RecordError(err)
				span.End()
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			switch tokenData.Type {
			case "ADMIN":
				c.Set(AdminIdParam, tokenData.Id)
				span.SetAttribute(AdminIdParam, tokenData.Id)
			case "CUSTOMER":
				c.Set(CustomerIdParam, tokenData.Id)
				span.SetAttribute(CustomerIdParam, tokenData.Id)
			default:
				span.RecordError(err)
				span.End()
				return errors.New("Invalid token type")
			}
			c.Set(tenant.Id, tokenData.TenantId)
			span.End()
			return next(c)
		}
	}
}

func OnlyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return requiredParam(AdminIdParam, next)
}

func OnlyCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return requiredParam(CustomerIdParam, next)
}

func requiredParam(name string, next echo.HandlerFunc) func(c echo.Context) error {
	return func(c echo.Context) error {
		param := c.Get(name)
		if param == nil {
			return errors.New("You dont have permission to access this resource")
		}
		return next(c)
	}
}
