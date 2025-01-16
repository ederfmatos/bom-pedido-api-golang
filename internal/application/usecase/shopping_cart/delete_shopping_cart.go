package shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"context"
)

type (
	DeleteShoppingCartUseCase struct {
		shoppingCartRepository repository.ShoppingCartRepository
	}
	DeleteShoppingCartInput struct {
		CustomerId string
	}
)

func NewDeleteShoppingCart(factory *factory.ApplicationFactory) *DeleteShoppingCartUseCase {
	return &DeleteShoppingCartUseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *DeleteShoppingCartUseCase) Execute(ctx context.Context, input DeleteShoppingCartInput) error {
	return useCase.shoppingCartRepository.DeleteByCustomerId(ctx, input.CustomerId)
}
