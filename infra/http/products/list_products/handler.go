package list_products

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/http/response"
	"bom-pedido-api/infra/tenant"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	return func(context echo.Context) error {
		output, err := factory.ProductQuery.List(context.Request().Context(), context.Get(tenant.Id).(string))
		return response.Ok(context, output, err)
	}
}
