package order

import (
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"time"
)

type ItemStatus string

const (
	OrderItemStatusOk        ItemStatus = "OK"
	OrderItemStatusCancelled ItemStatus = "CANCELLED"
)

type (
	Order struct {
		Id              string
		CustomerID      string
		PaymentMethod   enums.PaymentMethod
		PaymentMode     enums.PaymentMode
		DeliveryMode    enums.DeliveryMode
		CreatedAt       time.Time
		CreditCardToken string
		Change          float64
		Code            int32
		DeliveryTime    time.Time
		status          *Status
		Items           []Item
		History         []*StatusHistory
	}

	Item struct {
		Id          string
		ProductId   string
		Quantity    int
		Observation string
		Status      ItemStatus
		Price       float64
	}
)

func New(customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken string, change float64, deliveryTime time.Time) (*Order, error) {
	paymentMethod, deliveryMode, paymentMode, err := validateOrder(paymentMethodString, deliveryModeString, paymentModeString, creditCardToken, change)
	if err != nil {
		return nil, err
	}
	return &Order{
		Id:              value_object.NewID(),
		CustomerID:      customerID,
		PaymentMethod:   paymentMethod,
		PaymentMode:     paymentMode,
		DeliveryMode:    deliveryMode,
		CreatedAt:       time.Now(),
		CreditCardToken: creditCardToken,
		Change:          change,
		Code:            0,
		DeliveryTime:    deliveryTime,
		status:          AwaitingApproval,
		Items:           []Item{},
		History:         []*StatusHistory{},
	}, nil
}

func Restore(
	Id, customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken, status string,
	createdAt time.Time,
	change float64,
	code int32,
	deliveryTime time.Time,
	items []Item,
) (*Order, error) {
	paymentMethod, deliveryMode, paymentMode, err := validateOrder(paymentMethodString, deliveryModeString, paymentModeString, creditCardToken, change)
	if err != nil {
		return nil, err
	}
	orderStatus, err := ParseStatus(status)
	if err != nil {
		return nil, err
	}
	return &Order{
		Id:              Id,
		CustomerID:      customerID,
		PaymentMethod:   paymentMethod,
		PaymentMode:     paymentMode,
		DeliveryMode:    deliveryMode,
		CreatedAt:       createdAt,
		CreditCardToken: creditCardToken,
		Change:          change,
		Code:            code,
		DeliveryTime:    deliveryTime,
		status:          orderStatus,
		Items:           items,
		History:         []*StatusHistory{},
	}, nil
}

func validateOrder(paymentMethodString, deliveryModeString, paymentModeString, cardToken string, change float64) (enums.PaymentMethod, enums.DeliveryMode, enums.PaymentMode, error) {
	compositeError := errors.NewCompositeError()

	paymentMethod, err := enums.ParsePaymentMethod(paymentMethodString)
	compositeError.Append(err)

	deliveryMode, err := enums.ParseDeliveryMode(deliveryModeString)
	compositeError.Append(err)

	paymentMode, err := enums.ParsePaymentMode(paymentModeString)
	compositeError.Append(err)

	cardTokenIsRequired := paymentMethod.IsCreditCard() && paymentMode.IsInApp()
	if cardTokenIsRequired && cardToken == "" {
		compositeError.Append(errors.CardTokenIsRequiredError)
	}
	if change < 0 {
		compositeError.Append(errors.ChangeShouldBePositiveError)
	}
	return paymentMethod, deliveryMode, paymentMode, compositeError.AsError()
}

func (order *Order) AddProduct(product *product.Product, quantity int, observation string) *errors.DomainError {
	if product == nil {
		return errors.ProductNotFoundError
	}
	if product.IsUnAvailable() {
		return errors.ProductUnAvailableError
	}
	order.Items = append(order.Items, Item{
		Id:          value_object.NewID(),
		ProductId:   product.Id,
		Quantity:    quantity,
		Observation: observation,
		Price:       product.Price,
		Status:      OrderItemStatusOk,
	})
	return nil
}

func (order *Order) Approve(approvedAt time.Time, approvedBy string) error {
	approval, err := order.status.approve(approvedAt, approvedBy)
	if err != nil {
		return err
	}
	order.status = Approved
	order.History = append(order.History, approval)
	return nil
}

func (order *Order) GetStatus() string {
	return order.status.name
}
