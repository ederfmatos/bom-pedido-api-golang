package checkout_shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/shopping_cart"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/response"
	"github.com/labstack/echo/v4"
)

type checkoutShoppingCartRequest struct {
	PaymentMethod   string  `body:"paymentMethod" json:"paymentMethod,omitempty"`
	DeliveryMode    string  `body:"deliveryMode" json:"deliveryMode,omitempty"`
	PaymentMode     string  `body:"paymentMode" json:"paymentMode,omitempty"`
	AddressId       string  `body:"addressId" json:"addressId,omitempty"`
	Payback         float64 `body:"payback" json:"payback,omitempty"`
	CreditCardToken string  `body:"creditCardToken" json:"creditCardToken,omitempty"`
}

func Handle(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := shopping_cart.NewCheckoutShoppingCart(factory)
	return func(context echo.Context) error {
		var request checkoutShoppingCartRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := shopping_cart.CheckoutShoppingCartInput{
			CustomerId:      context.Get(middlewares.CustomerIdParam).(string),
			PaymentMethod:   request.PaymentMethod,
			DeliveryMode:    request.DeliveryMode,
			PaymentMode:     request.PaymentMode,
			AddressId:       request.AddressId,
			Payback:         request.Payback,
			CreditCardToken: request.CreditCardToken,
		}
		output, err := useCase.Execute(context.Request().Context(), input)
		return response.Ok(context, output, err)
	}
}
