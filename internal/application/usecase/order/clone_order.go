package order

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/shopping_cart"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	CloneOrderUseCase struct {
		orderRepository        repository.OrderRepository
		shoppingCartRepository repository.ShoppingCartRepository
	}
	CloneOrderInput struct {
		OrderId string
	}
)

func NewCloneOrder(factory *factory.ApplicationFactory) *CloneOrderUseCase {
	return &CloneOrderUseCase{
		orderRepository:        factory.OrderRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *CloneOrderUseCase) Execute(ctx context.Context, input CloneOrderInput) error {
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
