package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
	"time"
)

type DefaultOrderRepository struct {
	SqlConnection
}

func NewDefaultOrderRepository(sqlConnection SqlConnection) repository.OrderRepository {
	return &DefaultOrderRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	return repository.InTransaction(ctx, func(transaction SqlTransaction) error {
		err := transaction.Sql("INSERT INTO orders (id, customer_id, payment_method, payment_mode, delivery_mode, credit_card_token, `change`, delivery_time, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)").
			Values(order.Id, order.CustomerID, order.PaymentMethod, order.PaymentMode, order.DeliveryMode, order.CreditCardToken, order.Change, order.DeliveryTime, order.Status, order.CreatedAt).
			Update(ctx)
		if err != nil {
			return err
		}
		for _, item := range order.Items {
			err := transaction.Sql("INSERT INTO order_items (id, order_id, product_id, quantity, observation, price, status) VALUES (?, ?, ?, ?, ?, ?, ?)").
				Values(item.Id, order.Id, item.ProductId, item.Quantity, item.Observation, item.Price, item.Status).
				Update(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (repository *DefaultOrderRepository) FindById(ctx context.Context, id string) (*entity.Order, error) {
	var customerId, paymentMethod, paymentMode, deliveryMode, creditCardToken, status string
	var createdAt, deliveryTime time.Time
	var code int32
	var change float64
	found, err := repository.Sql(`
		SELECT code, customer_id, payment_method, payment_mode, delivery_mode, credit_card_token, change, delivery_time, status, created_at
		FROM orders
		WHERE id = ?
		LIMIT 1
	`).Values(id).
		FindOne(ctx, &code, &customerId, &paymentMethod, &paymentMode, &deliveryMode, &creditCardToken, &change, &deliveryTime, &status, &createdAt)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	var items []entity.OrderItem
	err = repository.Sql(`
		SELECT id, product_id, quantity, observation, price, status
		FROM order_items
		WHERE order_id = ?
	`).Values(id).List(ctx, func(getValues func(dest ...any) error) error {
		var item entity.OrderItem
		err := getValues(&item.Id, &item.ProductId, &item.Quantity, &item.Observation, &item.Price, &item.Status)
		if err != nil {
			return err
		}
		items = append(items, item)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entity.RestoreOrder(
		id,
		customerId,
		paymentMethod,
		paymentMode,
		deliveryMode,
		creditCardToken,
		status,
		createdAt,
		change,
		code,
		deliveryTime,
		[]entity.OrderItem{},
	)
}
