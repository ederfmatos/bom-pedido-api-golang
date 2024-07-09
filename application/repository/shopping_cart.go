package repository

import (
	"bom-pedido-api/domain/entity/shopping_cart"
	"context"
)

type ShoppingCartRepository interface {
	Upsert(ctx context.Context, shoppingCart *shopping_cart.ShoppingCart) error
	FindByCustomerId(ctx context.Context, id string) (*shopping_cart.ShoppingCart, error)
	DeleteByCustomerId(ctx context.Context, id string) error
}
