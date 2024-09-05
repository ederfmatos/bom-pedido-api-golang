package get_shopping_cart

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/projection"
	"bom-pedido-api/application/repository"
	"context"
)

type (
	UseCase struct {
		productRepository      repository.ProductRepository
		shoppingCartRepository repository.ShoppingCartRepository
	}
	Input struct {
		CustomerId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		productRepository:      factory.ProductRepository,
		shoppingCartRepository: factory.ShoppingCartRepository,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) (*projection.ShoppingCart, error) {
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
