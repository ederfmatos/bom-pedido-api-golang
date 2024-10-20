package add_item_to_shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/shopping_cart"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	UseCase struct {
		productRepository      repository.ProductRepository
		shoppingCartRepository repository.ShoppingCartRepository
	}
	Input struct {
		CustomerId  string
		ProductId   string
		Quantity    int
		Observation string
		TenantId    string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		productRepository:      factory.ProductRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) error {
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
		shoppingCart = shopping_cart.New(input.CustomerId, input.TenantId)
	}
	err = shoppingCart.AddItem(product, input.Quantity, input.Observation)
	if err != nil {
		return err
	}
	return useCase.shoppingCartRepository.Upsert(ctx, shoppingCart)
}
