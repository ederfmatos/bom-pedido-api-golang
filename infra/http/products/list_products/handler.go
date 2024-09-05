package list_products

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/projection"
	"bom-pedido-api/infra/http/response"
	"bom-pedido-api/infra/tenant"
	"github.com/labstack/echo/v4"
)

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	return func(context echo.Context) error {
		filter := projection.ProductListFilter{
			CurrentPage: 1,
			PageSize:    10,
			TenantId:    context.Get(tenant.Id).(string),
		}
		err := echo.QueryParamsBinder(context).
			Int32("pageSize", &filter.PageSize).
			Int32("currentPage", &filter.CurrentPage).
			BindError()
		if err != nil {
			return err
		}
		output, err := factory.ProductQuery.List(context.Request().Context(), filter)
		return response.Ok(context, output, err)
	}
}
