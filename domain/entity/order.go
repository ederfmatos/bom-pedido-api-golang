package entity

import (
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"time"
)

type Order struct {
	ID              string
	CustomerID      string
	PaymentMethod   enums.PaymentMethod
	PaymentMode     enums.PaymentMode
	DeliveryMode    enums.DeliveryMode
	CreatedAt       time.Time
	CreditCardToken string
	Change          float64
	Code            int32
	DeliveryTime    time.Time
	Items           []OrderItem
}

type OrderItem struct {
	Id          string
	ProductId   string
	Quantity    int
	Observation string
	Price       float64
}

func (order *Order) AddProduct(product *Product, quantity int, observation string) *errors.DomainError {
	if product == nil {
		return ProductNotFoundError
	}
	if product.IsUnAvailable() {
		return ProductUnAvailableError
	}
	order.Items = append(order.Items, OrderItem{
		Id:          value_object.NewID(),
		ProductId:   product.ID,
		Quantity:    quantity,
		Observation: observation,
		Price:       product.Price,
	})
	return nil
}

func NewOrder(customerID string, paymentMethod enums.PaymentMethod, paymentMode enums.PaymentMode, deliveryMode enums.DeliveryMode, creditCardToken string, change float64) *Order {
	return &Order{
		ID:              value_object.NewID(),
		CustomerID:      customerID,
		PaymentMethod:   paymentMethod,
		PaymentMode:     paymentMode,
		DeliveryMode:    deliveryMode,
		CreatedAt:       time.Now(),
		CreditCardToken: creditCardToken,
		Change:          change,
		Code:            0,
		DeliveryTime:    time.Now().Add(45 * time.Minute),
	}
}
