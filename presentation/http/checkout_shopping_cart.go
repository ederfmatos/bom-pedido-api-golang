package http

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase"
	"github.com/labstack/echo/v4"
)

type checkoutShoppingCartRequest struct {
	PaymentMethod   string  `body:"paymentMethod" json:"paymentMethod,omitempty"`
	DeliveryMode    string  `body:"deliveryMode" json:"deliveryMode,omitempty"`
	PaymentMode     string  `body:"paymentMode" json:"paymentMode,omitempty"`
	AddressId       string  `body:"addressId" json:"addressId,omitempty"`
	Change          float64 `body:"change" json:"change,omitempty"`
	CreditCardToken string  `body:"creditCardToken" json:"creditCardToken,omitempty"`
}

func HandleCheckoutShoppingCart(factory *factory.ApplicationFactory) func(context echo.Context) error {
	useCase := usecase.NewCheckoutShoppingCartUseCase(factory)
	return func(context echo.Context) error {
		var request checkoutShoppingCartRequest
		err := context.Bind(&request)
		if err != nil {
			return err
		}
		input := usecase.CheckoutShoppingCartInput{
			Context:         context.Request().Context(),
			CustomerId:      context.Get("customerId").(string),
			PaymentMethod:   request.PaymentMethod,
			DeliveryMode:    request.DeliveryMode,
			PaymentMode:     request.PaymentMode,
			AddressId:       request.AddressId,
			Change:          request.Change,
			CreditCardToken: request.CreditCardToken,
		}
		output, err := useCase.Execute(input)
		if err != nil {
			return err
		}
		return context.JSON(200, output)
	}
}
