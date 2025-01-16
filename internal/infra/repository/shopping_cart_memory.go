package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type ShoppingCartMemoryRepository struct {
	shoppingCarts map[string]*entity.ShoppingCart
}

func NewShoppingCartMemoryRepository() *ShoppingCartMemoryRepository {
	return &ShoppingCartMemoryRepository{shoppingCarts: make(map[string]*entity.ShoppingCart)}
}

func (r *ShoppingCartMemoryRepository) Upsert(_ context.Context, shoppingCart *entity.ShoppingCart) error {
	r.shoppingCarts[shoppingCart.CustomerId] = shoppingCart
	return nil
}

func (r *ShoppingCartMemoryRepository) DeleteByCustomerId(_ context.Context, id string) error {
	delete(r.shoppingCarts, id)
	return nil
}

func (r *ShoppingCartMemoryRepository) FindByCustomerId(_ context.Context, id string) (*entity.ShoppingCart, error) {
	return r.shoppingCarts[id], nil
}
