package usecase

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/domain/events"
	"context"
	"time"
)

type CheckoutShoppingCartUseCase struct {
	shoppingCartRepository repository.ShoppingCartRepository
	productRepository      repository.ProductRepository
	orderRepository        repository.OrderRepository
	eventEmitter           event.EventEmitter
}

func NewCheckoutShoppingCartUseCase(factory *factory.ApplicationFactory) *CheckoutShoppingCartUseCase {
	return &CheckoutShoppingCartUseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
		productRepository:      factory.ProductRepository,
		orderRepository:        factory.OrderRepository,
		eventEmitter:           factory.EventEmitter,
	}
}

type CheckoutShoppingCartInput struct {
	Context         context.Context
	CustomerId      string
	PaymentMethod   string
	DeliveryMode    string
	PaymentMode     string
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
	if shoppingCart == nil || shoppingCart.IsEmpty() {
		return nil, entity.ShoppingCartEmptyError
	}
	var productIds []string
	for _, item := range shoppingCart.GetItems() {
		productIds = append(productIds, item.ProductId)
	}
	products, err := useCase.productRepository.FindAllById(input.Context, productIds)
	if err != nil {
		return nil, err
	}
	// TODO: Change to get delivery time from merchant
	deliveryTimeInMinutes := 45 * time.Minute
	order, err := shoppingCart.Checkout(
		input.PaymentMethod,
		input.DeliveryMode,
		input.PaymentMode,
		input.CreditCardToken,
		input.Change,
		products,
		deliveryTimeInMinutes,
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
	return &CheckoutShoppingCartOutput{Id: order.Id}, nil
}
