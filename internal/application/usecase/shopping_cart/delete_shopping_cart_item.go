package shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	DeleteShoppingCartItemUseCase struct {
		shoppingCartRepository repository.ShoppingCartRepository
	}
	DeleteShoppingCartItemInput struct {
		CustomerId string
		ItemId     string
	}
)

func NewDeleteShoppingCartItem(factory *factory.ApplicationFactory) *DeleteShoppingCartItemUseCase {
	return &DeleteShoppingCartItemUseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *DeleteShoppingCartItemUseCase) Execute(ctx context.Context, input DeleteShoppingCartItemInput) error {
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(ctx, input.CustomerId)
	if err != nil {
		return err
	}
	if shoppingCart == nil {
		return errors.ShoppingCartEmptyError
	}
	shoppingCart.RemoveItem(input.ItemId)
	return useCase.shoppingCartRepository.Upsert(ctx, shoppingCart)
}
