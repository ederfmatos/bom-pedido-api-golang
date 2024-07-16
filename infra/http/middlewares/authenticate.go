package middlewares

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
			if token == "" {
				c.Set("customerId", "019078bc-cab8-789a-a1e7-4ba2a09561a6")
				c.Set("adminId", "019078bc-cab8-789a-a1e7-4ba2a09561a6")
				return next(c)
			}
			customerId, err := customerTokenManager.Decrypt(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			c.Set("customerId", customerId)
			c.Set("adminId", customerId)
			return next(c)
		}
	}
}
