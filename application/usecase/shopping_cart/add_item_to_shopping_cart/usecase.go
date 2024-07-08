package add_item_to_shopping_cart

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/shopping_cart"
	"bom-pedido-api/domain/errors"
	"context"
)

type (
	UseCase struct {
		productRepository      repository.ProductRepository
		shoppingCartRepository repository.ShoppingCartRepository
	}
	Input struct {
		Context     context.Context
		CustomerId  string
		ProductId   string
		Quantity    int
		Observation string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		productRepository:      factory.ProductRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *UseCase) Execute(input Input) error {
	product, err := useCase.productRepository.FindById(input.Context, input.ProductId)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.ProductNotFoundError
	}
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(input.Context, input.CustomerId)
	if err != nil {
		return err
	}
	if shoppingCart == nil {
		shoppingCart = shopping_cart.New(input.CustomerId)
	}
	err = shoppingCart.AddItem(product, input.Quantity, input.Observation)
	if err != nil {
		return err
	}
	return useCase.shoppingCartRepository.Upsert(input.Context, shoppingCart)
}
