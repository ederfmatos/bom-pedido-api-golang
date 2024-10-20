package delete_shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"context"
)

type (
	UseCase struct {
		shoppingCartRepository repository.ShoppingCartRepository
	}
	Input struct {
		CustomerId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) error {
	return useCase.shoppingCartRepository.DeleteByCustomerId(ctx, input.CustomerId)
}
