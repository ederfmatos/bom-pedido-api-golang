package repository

import (
	"bom-pedido-api/internal/domain/entity/shopping_cart"
	"context"
)

type ShoppingCartMemoryRepository struct {
	shoppingCarts map[string]*shopping_cart.ShoppingCart
}

func NewShoppingCartMemoryRepository() *ShoppingCartMemoryRepository {
	return &ShoppingCartMemoryRepository{shoppingCarts: make(map[string]*shopping_cart.ShoppingCart)}
}

func (r *ShoppingCartMemoryRepository) Upsert(_ context.Context, shoppingCart *shopping_cart.ShoppingCart) error {
	r.shoppingCarts[shoppingCart.CustomerId] = shoppingCart
	return nil
}

func (r *ShoppingCartMemoryRepository) DeleteByCustomerId(_ context.Context, id string) error {
	delete(r.shoppingCarts, id)
	return nil
}

func (r *ShoppingCartMemoryRepository) FindByCustomerId(_ context.Context, id string) (*shopping_cart.ShoppingCart, error) {
	return r.shoppingCarts[id], nil
}
