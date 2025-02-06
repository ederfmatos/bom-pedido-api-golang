package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
	"fmt"
)

type OrderStatusHistoryMongoRepository struct {
	collection mongo.Collection
}

func NewOrderStatusHistoryMongoRepository(database *mongo.Database) repository.OrderStatusHistoryRepository {
	return &OrderStatusHistoryMongoRepository{collection: database.ForCollection("order_status_history")}
}

func (r *OrderStatusHistoryMongoRepository) Create(ctx context.Context, history *entity.OrderStatusHistory) error {
	return r.collection.InsertOne(ctx, history)
}

func (r *OrderStatusHistoryMongoRepository) ListByOrderId(ctx context.Context, id string) ([]entity.OrderStatusHistory, error) {
	items := make([]entity.OrderStatusHistory, 0)
	cursor, err := r.collection.FindAllBy(ctx, map[string]interface{}{"orderId": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("decode history: %v", err)
	}
	return items, nil
}
