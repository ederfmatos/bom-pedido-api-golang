package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/shopping_cart"
	"context"
)

type ShoppingCartMemoryRepository struct {
	shoppingCarts map[string]*shopping_cart.ShoppingCart
}

func NewShoppingCartMemoryRepository() repository.ShoppingCartRepository {
	return &ShoppingCartMemoryRepository{shoppingCarts: make(map[string]*shopping_cart.ShoppingCart)}
}

func (repository *ShoppingCartMemoryRepository) Upsert(_ context.Context, shoppingCart *shopping_cart.ShoppingCart) error {
	repository.shoppingCarts[shoppingCart.CustomerId] = shoppingCart
	return nil
}

func (repository *ShoppingCartMemoryRepository) DeleteByCustomerId(ctx context.Context, id string) error {
	delete(repository.shoppingCarts, id)
	return nil
}

func (repository *ShoppingCartMemoryRepository) FindByCustomerId(_ context.Context, id string) (*shopping_cart.ShoppingCart, error) {
	return repository.shoppingCarts[id], nil
}
