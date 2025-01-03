package clone_order

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/shopping_cart"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	UseCase struct {
		orderRepository        repository.OrderRepository
		shoppingCartRepository repository.ShoppingCartRepository
	}
	Input struct {
		OrderId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		orderRepository:        factory.OrderRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) error {
	order, err := useCase.orderRepository.FindById(ctx, input.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.OrderNotFoundError
	}
	aShoppingCart := shopping_cart.CloneOrder(order)
	return useCase.shoppingCartRepository.Upsert(ctx, aShoppingCart)
}
