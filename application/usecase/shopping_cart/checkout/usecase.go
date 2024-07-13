package checkout

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/events"
	"context"
	"time"
)

type (
	UseCase struct {
		shoppingCartRepository repository.ShoppingCartRepository
		productRepository      repository.ProductRepository
		orderRepository        repository.OrderRepository
		eventEmitter           event.Emitter
	}
	Input struct {
		CustomerId      string
		PaymentMethod   string
		DeliveryMode    string
		PaymentMode     string
		AddressId       string // TODO
		Change          float64
		CreditCardToken string
	}
	Output struct {
		Id string `json:"id"`
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
		productRepository:      factory.ProductRepository,
		orderRepository:        factory.OrderRepository,
		eventEmitter:           factory.EventEmitter,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(ctx, input.CustomerId)
	if err != nil {
		return nil, err
	}
	if shoppingCart == nil || shoppingCart.IsEmpty() {
		return nil, errors.ShoppingCartEmptyError
	}
	var productIds []string
	for _, item := range shoppingCart.GetItems() {
		productIds = append(productIds, item.ProductId)
	}
	products, err := useCase.productRepository.FindAllById(ctx, productIds)
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
	err = useCase.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	err = useCase.eventEmitter.Emit(ctx, events.NewOrderCreatedEvent(order))
	if err != nil {
		return nil, err
	}
	return &Output{Id: order.Id}, nil
}
