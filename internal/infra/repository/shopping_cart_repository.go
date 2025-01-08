package repository

import (
	"bom-pedido-api/internal/domain/entity/shopping_cart"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type ShoppingCartMongoRepository struct {
	collection *mongo.Collection
}

func NewShoppingCartMongoRepository(database *mongo.Database) *ShoppingCartMongoRepository {
	return &ShoppingCartMongoRepository{collection: database.ForCollection("shopping_carts")}
}

func (r *ShoppingCartMongoRepository) Upsert(ctx context.Context, shoppingCart *shopping_cart.ShoppingCart) error {
	return r.collection.Upsert(ctx, shoppingCart.CustomerId, shoppingCart)
}

func (r *ShoppingCartMongoRepository) FindByCustomerId(ctx context.Context, id string) (*shopping_cart.ShoppingCart, error) {
	var shoppingCart shopping_cart.ShoppingCart
	err := r.collection.FindByID(ctx, id, &shoppingCart)
	if err != nil || shoppingCart.CustomerId == "" {
		return nil, err
	}
	return &shoppingCart, nil
}

func (r *ShoppingCartMongoRepository) DeleteByCustomerId(ctx context.Context, id string) error {
	return r.collection.DeleteByID(ctx, id)
}
