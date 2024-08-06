package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/entity/order/status"
	"context"
	"time"
)

const (
	sqlCreateOrder          = "INSERT INTO orders (id, customer_id, payment_method, payment_mode, delivery_mode, credit_card_token, change, delivery_time, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	sqlUpdateOrder          = "UPDATE orders SET status = $1 WHERE id = $2"
	sqlInsertOrderItem      = "INSERT INTO order_items (id, order_id, product_id, quantity, observation, price, status) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	sqlInsertOrderHistory   = "INSERT INTO order_history (order_id, changed_by, changed_at, status, data) VALUES ($1, $2, $3, $4, $5)"
	sqlFindOrderById        = "SELECT code, customer_id, payment_method, payment_mode, delivery_mode, credit_card_token, change, delivery_time, status, created_at FROM orders WHERE id = $1 LIMIT 1"
	sqlListItemsFromOrderId = "SELECT id, product_id, quantity, observation, price, status FROM order_items WHERE order_id = $1"
	sqlOrderHistory         = "SELECT changed_by, changed_at, status, data FROM order_history WHERE order_id = $1"
)

type (
	DefaultOrderRepository struct {
		SqlConnection
	}

	OrderEntity struct {
		Code            int32
		CustomerId      string
		PaymentMethod   string
		PaymentMode     string
		DeliveryMode    string
		CreditCardToken string
		Change          float64
		DeliveryTime    time.Time
		Status          string
		CreatedAt       time.Time
	}
)

func NewDefaultOrderRepository(sqlConnection SqlConnection) repository.OrderRepository {
	return &DefaultOrderRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultOrderRepository) Create(ctx context.Context, order *order.Order) error {
	return repository.InTransaction(ctx, func(transaction SqlTransaction, ctx context.Context) error {
		err := transaction.Sql(sqlCreateOrder).
			Values(order.Id, order.CustomerID, order.PaymentMethod, order.PaymentMode, order.DeliveryMode, order.CreditCardToken, order.Change, order.DeliveryTime, order.GetStatus(), order.CreatedAt).
			Update(ctx)

		if err != nil {
			return err
		}
		for _, item := range order.Items {
			err := transaction.Sql(sqlInsertOrderItem).
				Values(item.Id, order.Id, item.ProductId, item.Quantity, item.Observation, item.Price, item.Status).
				Update(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (repository *DefaultOrderRepository) FindById(ctx context.Context, id string) (*order.Order, error) {
	var entity OrderEntity
	found, err := repository.Sql(sqlFindOrderById).Values(id).
		FindOne(ctx, &entity.Code, &entity.CustomerId, &entity.PaymentMethod, &entity.PaymentMode, &entity.DeliveryMode, &entity.CreditCardToken, &entity.Change, &entity.DeliveryTime, &entity.Status, &entity.CreatedAt)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	items, err := repository.getOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	history, err := repository.getOrderHistory(ctx, id)
	if err != nil {
		return nil, err
	}
	return order.Restore(
		id,
		entity.CustomerId,
		entity.PaymentMethod,
		entity.PaymentMode,
		entity.DeliveryMode,
		entity.CreditCardToken,
		entity.Status,
		entity.CreatedAt,
		entity.Change,
		entity.Code,
		entity.DeliveryTime,
		items,
		history,
	)
}

func (repository *DefaultOrderRepository) getOrderItems(ctx context.Context, id string) ([]order.Item, error) {
	items := make([]order.Item, 0)
	err := repository.Sql(sqlListItemsFromOrderId).Values(id).List(ctx, func(getValues func(dest ...any) error) error {
		var item order.Item
		err := getValues(&item.Id, &item.ProductId, &item.Quantity, &item.Observation, &item.Price, &item.Status)
		if err != nil {
			return err
		}
		items = append(items, item)
		return nil
	})
	return items, err
}

func (repository *DefaultOrderRepository) getOrderHistory(ctx context.Context, id string) ([]status.History, error) {
	history := make([]status.History, 0)
	err := repository.Sql(sqlOrderHistory).Values(id).List(ctx, func(getValues func(dest ...any) error) error {
		var item status.History
		err := getValues(&item.ChangedBy, &item.Time, &item.Status, &item.Data)
		if err != nil {
			return err
		}
		item.Time = parseTime(item.Time)
		history = append(history, item)
		return nil
	})
	return history, err
}

func (repository *DefaultOrderRepository) Update(ctx context.Context, order *order.Order) error {
	return repository.InTransaction(ctx, func(transaction SqlTransaction, ctx context.Context) error {
		err := transaction.Sql(sqlUpdateOrder).
			Values(order.GetStatus(), order.Id).
			Update(ctx)
		if err != nil {
			return err
		}
		history := &order.History[len(order.History)-1]
		history.Time = parseTime(history.Time)
		return transaction.Sql(sqlInsertOrderHistory).
			Values(order.Id, history.ChangedBy, history.Time, history.Status, history.Data).
			Update(ctx)
	})
}

func parseTime(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(), (value.Nanosecond()/1000)*1000, time.Local)
}
