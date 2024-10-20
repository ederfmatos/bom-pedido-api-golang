package delete_shopping_cart_item

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	UseCase struct {
		shoppingCartRepository repository.ShoppingCartRepository
	}
	Input struct {
		CustomerId string
		ItemId     string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) error {
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
