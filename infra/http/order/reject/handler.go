package reject

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/reject_order"
	"bom-pedido-api/infra/http/response"
	"github.com/labstack/echo/v4"
)

type rejectOrderRequest struct {
	Reason string `body:"reason" json:"reason,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := reject_order.New(factory)
	return func(context echo.Context) error {
		var request rejectOrderRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := reject_order.Input{
			OrderId:    context.Param("id"),
			RejectedBy: context.Get("adminId").(string),
			Reason:     request.Reason,
		}
		err = useCase.Execute(context.Request().Context(), input)
		return response.NoContent(context, err)
	}
}
