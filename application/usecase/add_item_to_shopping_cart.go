package usecase

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type AddItemToShoppingCartInput struct {
	Context     context.Context
	CustomerId  string
	ProductId   string
	Quantity    int
	Observation string
}

type AddItemToShoppingCartUseCase struct {
	productRepository      repository.ProductRepository
	shoppingCartRepository repository.ShoppingCartRepository
}

func NewAddItemToShoppingCartUseCase(factory *factory.ApplicationFactory) *AddItemToShoppingCartUseCase {
	return &AddItemToShoppingCartUseCase{
		productRepository:      factory.ProductRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *AddItemToShoppingCartUseCase) Execute(input AddItemToShoppingCartInput) error {
	product, err := useCase.productRepository.FindById(input.Context, input.ProductId)
	if err != nil {
		return err
	}
	if product == nil {
		return entity.ProductNotFoundError
	}
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(input.Context, input.CustomerId)
	if err != nil {
		return err
	}
	if shoppingCart == nil {
		shoppingCart = entity.NewShoppingCart(input.CustomerId)
	}
	err = shoppingCart.AddItem(product, input.Quantity, input.Observation)
	if err != nil {
		return err
	}
	return useCase.shoppingCartRepository.Upsert(input.Context, shoppingCart)
}
