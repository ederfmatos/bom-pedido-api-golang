package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/order"
	"context"
)

const (
	sqlInsertOrderStatusHistory = "INSERT INTO order_history (order_id, changed_by, changed_at, status, data) VALUES ($1, $2, $3, $4, $5)"
	sqlOrderHistory             = "SELECT changed_by, changed_at, status, data FROM order_history WHERE order_id = $1"
)

type (
	DefaultOrderStatusHistoryRepository struct {
		SqlConnection
	}
)

func NewDefaultOrderStatusHistoryRepository(sqlConnection SqlConnection) repository.OrderStatusHistoryRepository {
	return &DefaultOrderStatusHistoryRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultOrderStatusHistoryRepository) Create(ctx context.Context, history *order.StatusHistory) error {
	return repository.Sql(sqlInsertOrderStatusHistory).
		Values(history.OrderId, history.ChangedBy, history.Time, history.Status, history.Data).
		Update(ctx)
}

func (repository *DefaultOrderStatusHistoryRepository) ListByOrderId(ctx context.Context, id string) ([]order.StatusHistory, error) {
	history := make([]order.StatusHistory, 0)
	err := repository.Sql(sqlOrderHistory).
		Values(id).
		List(ctx, func(getValues func(dest ...any) error) error {
			var item order.StatusHistory
			err := getValues(&item.ChangedBy, &item.Time, &item.Status, &item.Data)
			if err != nil {
				return err
			}
			item.OrderId = id
			history = append(history, item)
			return nil
		})
	return history, err
}
