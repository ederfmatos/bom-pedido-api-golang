package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/order"
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	sqlCreateOrder          = "INSERT INTO orders (id, customer_id, payment_method, payment_mode, delivery_mode, credit_card_token, payback, amount, delivery_time, status, created_at, merchant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	sqlUpdateOrder          = "UPDATE orders SET status = $1 WHERE id = $2"
	sqlInsertOrderItem      = "INSERT INTO order_items (id, order_id, product_id, quantity, observation, price, status) VALUES "
	sqlFindOrderById        = "SELECT code, customer_id, payment_method, payment_mode, delivery_mode, credit_card_token, payback, amount, delivery_time, status, created_at, merchant_id FROM orders WHERE id = $1 LIMIT 1"
	sqlListItemsFromOrderId = "SELECT id, product_id, quantity, observation, price, status FROM order_items WHERE order_id = $1"
	orderItemsFieldSize     = 7
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
		Payback         float64
		Amount          float64
		DeliveryTime    time.Time
		Status          string
		CreatedAt       time.Time
		MerchantId      string
	}
)

func NewDefaultOrderRepository(sqlConnection SqlConnection) repository.OrderRepository {
	return &DefaultOrderRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultOrderRepository) Create(ctx context.Context, order *order.Order) error {
	return repository.InTransaction(ctx, func(transaction SqlTransaction, ctx context.Context) error {
		err := transaction.Sql(sqlCreateOrder).
			Values(order.Id, order.CustomerID, order.PaymentMethod, order.PaymentMode, order.DeliveryMode, order.CreditCardToken, order.Payback, order.Amount, order.DeliveryTime, order.GetStatus(), order.CreatedAt, order.MerchantId).
			Update(ctx)
		if err != nil {
			return err
		}
		itemsSql, values := repository.orderItemSql(order)
		return transaction.Sql(itemsSql).
			Values(values...).
			Update(ctx)
	})
}

func (repository *DefaultOrderRepository) orderItemSql(order *order.Order) (string, []interface{}) {
	var valuesSql strings.Builder
	valuesSql.WriteString(sqlInsertOrderItem)
	itemsSize := len(order.Items)
	values := make([]interface{}, itemsSize*orderItemsFieldSize)
	for i, item := range order.Items {
		valuesSql.WriteString("(")
		for j := 1; j <= orderItemsFieldSize; j++ {
			value := i*orderItemsFieldSize + j
			valuesSql.WriteString(fmt.Sprintf("$%d", value))
			if j < orderItemsFieldSize {
				valuesSql.WriteString(",")
			}
		}
		valuesSql.WriteString(")")
		if i < itemsSize-1 {
			valuesSql.WriteString(",")
		}
		fieldNumber := i * orderItemsFieldSize
		values[fieldNumber] = item.Id
		values[fieldNumber+1] = order.Id
		values[fieldNumber+2] = item.ProductId
		values[fieldNumber+3] = item.Quantity
		values[fieldNumber+4] = item.Observation
		values[fieldNumber+5] = item.Price
		values[fieldNumber+6] = item.Status
	}
	return valuesSql.String(), values
}

func (repository *DefaultOrderRepository) FindById(ctx context.Context, id string) (*order.Order, error) {
	var entity OrderEntity
	found, err := repository.Sql(sqlFindOrderById).Values(id).
		FindOne(ctx, &entity.Code, &entity.CustomerId, &entity.PaymentMethod, &entity.PaymentMode, &entity.DeliveryMode, &entity.CreditCardToken, &entity.Payback, &entity.Amount, &entity.DeliveryTime, &entity.Status, &entity.CreatedAt, &entity.MerchantId)
	if err != nil || !found {
		return nil, err
	}
	items, err := repository.getOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	return order.Restore(id, entity.CustomerId, entity.PaymentMethod, entity.PaymentMode, entity.DeliveryMode, entity.CreditCardToken, entity.Status, entity.CreatedAt, entity.Payback, entity.Amount, entity.Code, entity.DeliveryTime, items, entity.MerchantId)
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

func (repository *DefaultOrderRepository) Update(ctx context.Context, order *order.Order) error {
	return repository.Sql(sqlUpdateOrder).
		Values(order.GetStatus(), order.Id).
		Update(ctx)
}

func parseTime(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(), (value.Nanosecond()/1000)*1000, time.Local)
}
