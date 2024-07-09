package delete_shopping_cart

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"context"
)

type (
	UseCase struct {
		shoppingCartRepository repository.ShoppingCartRepository
	}
	Input struct {
		Context    context.Context
		CustomerId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *UseCase) Execute(input Input) error {
	return useCase.shoppingCartRepository.DeleteByCustomerId(input.Context, input.CustomerId)
}
