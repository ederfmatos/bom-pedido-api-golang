package usecase

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/events"
	"context"
)

type CheckoutShoppingCartUseCase struct {
	shoppingCartRepository repository.ShoppingCartRepository
	productRepository      repository.ProductRepository
	orderRepository        repository.OrderRepository
	eventEmitter           event.EventEmitter
}

type CheckoutShoppingCartInput struct {
	Context         context.Context
	CustomerId      string
	PaymentMethod   enums.PaymentMethod
	DeliveryMode    enums.DeliveryMode
	PaymentMode     enums.PaymentMode
	AddressId       string // TODO
	Change          float64
	CreditCardToken string
}

type CheckoutShoppingCartOutput struct {
	Id string `json:"id"`
}

func (useCase CheckoutShoppingCartUseCase) Execute(input CheckoutShoppingCartInput) (*CheckoutShoppingCartOutput, error) {
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(input.Context, input.CustomerId)
	if err != nil {
		return nil, err
	}
	if shoppingCart == nil {
		return nil, entity.ShoppingCartNotFoundError
	}
	var productIds []string
	for _, item := range shoppingCart.GetItems() {
		productIds = append(productIds, item.ProductId)
	}
	products, err := useCase.productRepository.FindAllById(input.Context, productIds)
	if err != nil {
		return nil, err
	}
	order, err := shoppingCart.Checkout(
		input.PaymentMethod,
		input.DeliveryMode,
		input.PaymentMode,
		input.Change,
		input.CreditCardToken,
		products,
	)
	if err != nil {
		return nil, err
	}
	err = useCase.orderRepository.Create(input.Context, order)
	if err != nil {
		return nil, err
	}
	err = useCase.eventEmitter.Emit(input.Context, events.NewOrderCreatedEvent(order))
	if err != nil {
		return nil, err
	}
	return &CheckoutShoppingCartOutput{Id: order.ID}, nil
}
