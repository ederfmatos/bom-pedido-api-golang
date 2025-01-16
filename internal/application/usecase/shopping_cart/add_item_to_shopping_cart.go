package shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	AddItemToShoppingCartUseCase struct {
		productRepository      repository.ProductRepository
		shoppingCartRepository repository.ShoppingCartRepository
	}
	AddItemToShoppingCartInput struct {
		CustomerId  string
		ProductId   string
		Quantity    int
		Observation string
		TenantId    string
	}
)

func NewAddItemToShoppingCart(factory *factory.ApplicationFactory) *AddItemToShoppingCartUseCase {
	return &AddItemToShoppingCartUseCase{
		productRepository:      factory.ProductRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *AddItemToShoppingCartUseCase) Execute(ctx context.Context, input AddItemToShoppingCartInput) error {
	product, err := useCase.productRepository.FindById(ctx, input.ProductId)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.ProductNotFoundError
	}
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(ctx, input.CustomerId)
	if err != nil {
		return err
	}
	if shoppingCart == nil {
		shoppingCart = entity.NewShoppingCart(input.CustomerId, input.TenantId)
	}
	err = shoppingCart.AddItem(product, input.Quantity, input.Observation)
	if err != nil {
		return err
	}
	return useCase.shoppingCartRepository.Upsert(ctx, shoppingCart)
}
