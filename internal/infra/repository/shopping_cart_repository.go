package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type ShoppingCartMongoRepository struct {
	collection *mongo.Collection
}

func NewShoppingCartMongoRepository(database *mongo.Database) *ShoppingCartMongoRepository {
	return &ShoppingCartMongoRepository{collection: database.ForCollection("shopping_carts")}
}

func (r *ShoppingCartMongoRepository) Upsert(ctx context.Context, shoppingCart *entity.ShoppingCart) error {
	return r.collection.Upsert(ctx, shoppingCart.CustomerId, shoppingCart)
}

func (r *ShoppingCartMongoRepository) FindByCustomerId(ctx context.Context, id string) (*entity.ShoppingCart, error) {
	var shoppingCart entity.ShoppingCart
	err := r.collection.FindByID(ctx, id, &shoppingCart)
	if err != nil || shoppingCart.CustomerId == "" {
		return nil, err
	}
	return &shoppingCart, nil
}

func (r *ShoppingCartMongoRepository) DeleteByCustomerId(ctx context.Context, id string) error {
	return r.collection.DeleteByID(ctx, id)
}
