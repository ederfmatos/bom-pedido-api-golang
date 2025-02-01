package handler

import (
	usecase "bom-pedido-api/internal/application/usecase/shopping_cart"
	"bom-pedido-api/pkg/http"
	"fmt"
)

type (
	ShoppingCartHandler struct {
		addItemToShoppingCartUseCase  *usecase.AddItemToShoppingCartUseCase
		checkoutShoppingCartUseCase   *usecase.CheckoutShoppingCartUseCase
		getShoppingCartUseCase        *usecase.GetShoppingCartUseCase
		deleteShoppingCartUseCase     *usecase.DeleteShoppingCartUseCase
		deleteShoppingCartItemUseCase *usecase.DeleteShoppingCartItemUseCase
	}

	addItemToShoppingCartRequest struct {
		ProductId   string `body:"productId" json:"productId,omitempty"`
		Quantity    int    `body:"quantity" json:"quantity,omitempty"`
		Observation string `body:"observation" json:"observation,omitempty"`
	}

	checkoutShoppingCartRequest struct {
		PaymentMethod   string  `body:"paymentMethod" json:"paymentMethod,omitempty"`
		DeliveryMode    string  `body:"deliveryMode" json:"deliveryMode,omitempty"`
		PaymentMode     string  `body:"paymentMode" json:"paymentMode,omitempty"`
		AddressId       string  `body:"addressId" json:"addressId,omitempty"`
		Payback         float64 `body:"payback" json:"payback,omitempty"`
		CreditCardToken string  `body:"creditCardToken" json:"creditCardToken,omitempty"`
	}
)

func NewShoppingCartHandler(addItemToShoppingCartUseCase *usecase.AddItemToShoppingCartUseCase, checkoutShoppingCartUseCase *usecase.CheckoutShoppingCartUseCase, getShoppingCartUseCase *usecase.GetShoppingCartUseCase, deleteShoppingCartUseCase *usecase.DeleteShoppingCartUseCase, deleteShoppingCartItemUseCase *usecase.DeleteShoppingCartItemUseCase) *ShoppingCartHandler {
	return &ShoppingCartHandler{addItemToShoppingCartUseCase: addItemToShoppingCartUseCase, checkoutShoppingCartUseCase: checkoutShoppingCartUseCase, getShoppingCartUseCase: getShoppingCartUseCase, deleteShoppingCartUseCase: deleteShoppingCartUseCase, deleteShoppingCartItemUseCase: deleteShoppingCartItemUseCase}
}

func (h ShoppingCartHandler) GetShoppingCart(request http.Request, response http.Response) error {
	input := usecase.GetShoppingCartInput{CustomerId: request.AuthenticatedUser()}
	shoppingCart, err := h.getShoppingCartUseCase.Execute(request.Context(), input)
	if err != nil {
		return fmt.Errorf("get shopping cart: %w", err)
	}

	return response.OK(shoppingCart)
}

func (h ShoppingCartHandler) AddShoppingCartItem(request http.Request, response http.Response) error {
	var body addItemToShoppingCartRequest
	if err := request.Bind(&body); err != nil {
		return err
	}

	input := usecase.AddItemToShoppingCartInput{
		CustomerId:  request.AuthenticatedUser(),
		ProductId:   body.ProductId,
		Quantity:    body.Quantity,
		Observation: body.Observation,
		TenantId:    request.TenantID(),
	}
	if err := h.addItemToShoppingCartUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("add item to shoppin cart: %w", err)
	}

	return response.NoContent()
}

func (h ShoppingCartHandler) Checkout(request http.Request, response http.Response) error {
	var body checkoutShoppingCartRequest
	if err := request.Bind(&body); err != nil {
		return err
	}

	input := usecase.CheckoutShoppingCartInput{
		CustomerId:      request.AuthenticatedUser(),
		PaymentMethod:   body.PaymentMethod,
		DeliveryMode:    body.DeliveryMode,
		PaymentMode:     body.PaymentMode,
		AddressId:       body.AddressId,
		Payback:         body.Payback,
		CreditCardToken: body.CreditCardToken,
	}
	output, err := h.checkoutShoppingCartUseCase.Execute(request.Context(), input)
	if err != nil {
		return fmt.Errorf("checkout shopping cart: %w", err)
	}

	return response.OK(output)
}

func (h ShoppingCartHandler) DeleteShoppingCart(request http.Request, response http.Response) error {
	input := usecase.DeleteShoppingCartInput{CustomerId: request.AuthenticatedUser()}
	if err := h.deleteShoppingCartUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("delete shopping cart: %w", err)
	}

	return response.NoContent()
}

func (h ShoppingCartHandler) DeleteShoppingCartItem(request http.Request, response http.Response) error {
	input := usecase.DeleteShoppingCartItemInput{
		CustomerId: request.AuthenticatedUser(),
		ItemId:     request.PathParam("id"),
	}
	if err := h.deleteShoppingCartItemUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("delete shopping cart: %w", err)
	}

	return response.NoContent()
}
