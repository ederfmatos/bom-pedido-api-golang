package shopping_cart

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/projection"
	"bom-pedido-api/internal/application/repository"
	"context"
)

type (
	GetShoppingCartUseCase struct {
		productRepository      repository.ProductRepository
		shoppingCartRepository repository.ShoppingCartRepository
	}
	GetShoppingCartInput struct {
		CustomerId string
	}
)

func NewGetShoppingCart(factory *factory.ApplicationFactory) *GetShoppingCartUseCase {
	return &GetShoppingCartUseCase{
		productRepository:      factory.ProductRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *GetShoppingCartUseCase) Execute(ctx context.Context, input GetShoppingCartInput) (*projection.ShoppingCart, error) {
	shoppingCart, err := useCase.shoppingCartRepository.FindByCustomerId(ctx, input.CustomerId)
	if err != nil {
		return nil, err
	}
	if shoppingCart == nil {
		return &projection.ShoppingCart{Amount: 0.0, Items: make([]projection.ShoppingCartItem, 0)}, nil
	}
	output := &projection.ShoppingCart{
		Amount: shoppingCart.GetPrice(),
		Items:  make([]projection.ShoppingCartItem, len(shoppingCart.Items)),
	}
	index := 0
	for _, item := range shoppingCart.Items {
		output.Items[index] = projection.ShoppingCartItem{
			Id:          item.Id,
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Price:       item.Price,
			TotalPrice:  item.GetTotalPrice(),
			Quantity:    item.Quantity,
		}
		index++
	}
	return output, nil
}
