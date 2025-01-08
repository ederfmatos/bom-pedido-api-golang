package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/order"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type OrderMongoRepository struct {
	collection *mongo.Collection
}

func NewOrderMongoRepository(database *mongo.Database) repository.OrderRepository {
	return &OrderMongoRepository{collection: database.ForCollection("orders")}
}

func (r *OrderMongoRepository) Create(ctx context.Context, order *order.Order) error {
	return r.collection.InsertOne(ctx, order)
}

func (r *OrderMongoRepository) FindById(ctx context.Context, id string) (*order.Order, error) {
	var anOrder order.Order
	err := r.collection.FindByID(ctx, id, &anOrder)
	if err != nil || anOrder.Id == "" {
		return nil, err
	}
	return &anOrder, nil
}

func (r *OrderMongoRepository) Update(ctx context.Context, order *order.Order) error {
	return r.collection.UpdateByID(ctx, order.Id, order)
}
