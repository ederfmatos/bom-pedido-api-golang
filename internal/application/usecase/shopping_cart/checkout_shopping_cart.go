package shopping_cart

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/errors"
	"context"
	"time"
)

type (
	CheckoutShoppingCartUseCase struct {
		shoppingCartRepository repository.ShoppingCartRepository
		merchantRepository     repository.MerchantRepository
		productRepository      repository.ProductRepository
		orderRepository        repository.OrderRepository
		eventEmitter           event.Emitter
	}
	CheckoutShoppingCartInput struct {
		CustomerId      string
		PaymentMethod   string
		DeliveryMode    string
		PaymentMode     string
		AddressId       string // TODO
		Payback         float64
		CreditCardToken string
	}
	CheckoutShoppingCartOutput struct {
		Id string `json:"id"`
	}
)

func NewCheckoutShoppingCart(factory *factory.ApplicationFactory) *CheckoutShoppingCartUseCase {
	return &CheckoutShoppingCartUseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
		merchantRepository:     factory.MerchantRepository,
		productRepository:      factory.ProductRepository,
		orderRepository:        factory.OrderRepository,
		eventEmitter:           factory.EventEmitter,
	}
}

func (useCase *CheckoutShoppingCartUseCase) Execute(ctx context.Context, input CheckoutShoppingCartInput) (*CheckoutShoppingCartOutput, error) {
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(ctx, input.CustomerId)
	if err != nil {
		return nil, err
	}
	if shoppingCart == nil || shoppingCart.IsEmpty() {
		return nil, errors.ShoppingCartEmptyError
	}
	merchant, err := useCase.merchantRepository.FindByTenantId(ctx, shoppingCart.TenantId)
	if err != nil || merchant == nil {
		return nil, err
	}
	var productIds []string
	for _, item := range shoppingCart.Items {
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
		input.Payback,
		products,
		deliveryTimeInMinutes,
		merchant.Id,
	)
	if err != nil {
		return nil, err
	}
	err = useCase.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	err = useCase.eventEmitter.Emit(ctx, event.NewOrderCreatedEvent(order))
	if err != nil {
		return nil, err
	}
	return &CheckoutShoppingCartOutput{Id: order.Id}, nil
}
