package repository

import (
	"bom-pedido-api/domain/entity/shopping_cart"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShoppingCartMongoRepository struct {
	collection *mongo.Collection
}

func NewShoppingCartMongoRepository(database *mongo.Database) *ShoppingCartMongoRepository {
	return &ShoppingCartMongoRepository{collection: database.Collection("shopping_carts")}
}

func (repository *ShoppingCartMongoRepository) Upsert(ctx context.Context, shoppingCart *shopping_cart.ShoppingCart) error {
	update := bson.M{"$set": shoppingCart}
	updateOptions := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: shoppingCart.CustomerId}}
	_, err := repository.collection.UpdateOne(ctx, filter, update, updateOptions)
	return err
}

func (repository *ShoppingCartMongoRepository) FindByCustomerId(ctx context.Context, id string) (*shopping_cart.ShoppingCart, error) {
	result := repository.collection.FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	shoppingCart := &shopping_cart.ShoppingCart{}
	err := result.Decode(shoppingCart)
	if err != nil {
		return nil, err
	}
	if shoppingCart.CustomerId == "" {
		return nil, nil
	}
	return shoppingCart, nil
}

func (repository *ShoppingCartMongoRepository) DeleteByCustomerId(ctx context.Context, id string) error {
	_, err := repository.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
