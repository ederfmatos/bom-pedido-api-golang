package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/shopping_cart"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

type ShoppingCartRedisRepository struct {
	*redis.Client
}

func NewShoppingCartRedisRepository(client *redis.Client) repository.ShoppingCartRepository {
	return &ShoppingCartRedisRepository{client}
}

func (repository *ShoppingCartRedisRepository) Upsert(context context.Context, shoppingCart *shopping_cart.ShoppingCart) error {
	value, err := json.Marshal(shoppingCart)
	if err != nil {
		return err
	}
	return repository.Set(context, "SHOPPING_CART::"+shoppingCart.CustomerId, value, 0).Err()
}

func (repository *ShoppingCartRedisRepository) FindByCustomerId(context context.Context, id string) (*shopping_cart.ShoppingCart, error) {
	value, err := repository.Get(context, "SHOPPING_CART::"+id).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var shoppingCart shopping_cart.ShoppingCart
	err = json.Unmarshal([]byte(value), &shoppingCart)
	return &shoppingCart, err
}
