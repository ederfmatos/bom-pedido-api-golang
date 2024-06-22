package middleware

import (
	"bom-pedido-api/application/factory"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthenticateMiddleware(factory *factory.ApplicationFactory) echo.MiddlewareFunc {
	customerTokenManager := factory.TokenFactory.CustomerTokenManager
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token != "" {
				return next(c)
			}
			customerId, err := customerTokenManager.Decrypt(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			c.Set("customerId", customerId)
			return next(c)
		}
	}
}
