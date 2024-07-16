package list_products

import (
	"bom-pedido-api/application/factory"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	return func(context echo.Context) error {
		output, err := factory.ProductQuery.List(context.Request().Context())
		if err != nil {
			return err
		}
		return context.JSON(200, output)
	}
}
