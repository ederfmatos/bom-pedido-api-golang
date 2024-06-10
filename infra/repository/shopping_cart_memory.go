package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type ShoppingCartMemoryRepository struct {
	shoppingCarts map[string]*entity.ShoppingCart
}

func NewShoppingCartMemoryRepository() repository.ShoppingCartRepository {
	return &ShoppingCartMemoryRepository{shoppingCarts: make(map[string]*entity.ShoppingCart)}
}

func (repository *ShoppingCartMemoryRepository) Upsert(_ context.Context, shoppingCart *entity.ShoppingCart) error {
	repository.shoppingCarts[shoppingCart.CustomerId] = shoppingCart
	return nil
}

func (repository *ShoppingCartMemoryRepository) FindByCustomerId(_ context.Context, id string) (*entity.ShoppingCart, error) {
	return repository.shoppingCarts[id], nil
}
