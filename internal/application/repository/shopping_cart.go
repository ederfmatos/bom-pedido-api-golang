package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type ShoppingCartRepository interface {
	Upsert(ctx context.Context, shoppingCart *entity.ShoppingCart) error
	FindByCustomerId(ctx context.Context, id string) (*entity.ShoppingCart, error)
	DeleteByCustomerId(ctx context.Context, id string) error
}
