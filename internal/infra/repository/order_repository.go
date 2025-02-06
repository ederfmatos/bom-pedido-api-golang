package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type OrderMongoRepository struct {
	collection mongo.Collection
}

func NewOrderMongoRepository(database *mongo.Database) repository.OrderRepository {
	return &OrderMongoRepository{collection: database.ForCollection("orders")}
}

func (r *OrderMongoRepository) Create(ctx context.Context, order *entity.Order) error {
	return r.collection.InsertOne(ctx, order)
}

func (r *OrderMongoRepository) FindById(ctx context.Context, id string) (*entity.Order, error) {
	var order entity.Order
	err := r.collection.FindByID(ctx, id, &order)
	if err != nil || order.Id == "" {
		return nil, err
	}
	return &order, nil
}

func (r *OrderMongoRepository) Update(ctx context.Context, order *entity.Order) error {
	return r.collection.UpdateByID(ctx, order.Id, order)
}
