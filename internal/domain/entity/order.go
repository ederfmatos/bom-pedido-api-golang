package entity

import (
	"bom-pedido-api/internal/domain/entity/status"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"time"
)

const (
	OrderItemStatusOk OrderItemStatus = "OK"
)

type (
	OrderItemStatus string

	Order struct {
		Id              string              `bson:"id"`
		CustomerID      string              `bson:"customerID"`
		PaymentMethod   enums.PaymentMethod `bson:"paymentMethod"`
		PaymentMode     enums.PaymentMode   `bson:"paymentMode"`
		DeliveryMode    enums.DeliveryMode  `bson:"deliveryMode"`
		CreatedAt       time.Time           `bson:"createdAt"`
		CreditCardToken string              `bson:"creditCardToken"`
		Payback         float64             `bson:"payback"`
		Code            int32               `bson:"code"`
		DeliveryTime    time.Time           `bson:"deliveryTime"`
		Status          string              `bson:"status"`
		Items           []OrderItem         `bson:"items"`
		MerchantId      string              `bson:"merchantId"`
		Amount          float64             `bson:"amount"`
	}

	OrderItem struct {
		Id          string          `bson:"id"`
		ProductId   string          `bson:"productId"`
		Quantity    int             `bson:"quantity"`
		Observation string          `bson:"observation"`
		Status      OrderItemStatus `bson:"status"`
		Price       float64         `bson:"price"`
	}
)

func NewOrder(customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken string, payback, amount float64, deliveryTime time.Time, merchantId string) (*Order, error) {
	paymentMethod, deliveryMode, paymentMode, err := validateOrder(paymentMethodString, deliveryModeString, paymentModeString, creditCardToken, payback)
	if err != nil {
		return nil, err
	}
	state := status.AwaitingApprovalStatus
	if paymentMode.IsInApp() {
		state = status.AwaitingPaymentStatus
	}
	return &Order{
		Id:              value_object.NewID(),
		CustomerID:      customerID,
		PaymentMethod:   paymentMethod,
		PaymentMode:     paymentMode,
		DeliveryMode:    deliveryMode,
		CreatedAt:       time.Now().UTC(),
		CreditCardToken: creditCardToken,
		Payback:         payback,
		Code:            0,
		DeliveryTime:    deliveryTime.UTC(),
		Items:           make([]OrderItem, 0),
		MerchantId:      merchantId,
		Amount:          amount,
		Status:          state.Name(),
	}, nil
}

// TODO: Remover. Nos testes trocar por fixture
func RestoreOrder(Id, customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken, orderStatusString string, createdAt time.Time, payback, amount float64, code int32, deliveryTime time.Time, items []OrderItem, merchantId string) (*Order, error) {
	paymentMethod, deliveryMode, paymentMode, err := validateOrder(paymentMethodString, deliveryModeString, paymentModeString, creditCardToken, payback)
	if err != nil {
		return nil, err
	}
	orderStatus, err := status.Parse(orderStatusString)
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
		Payback:         payback,
		Code:            code,
		DeliveryTime:    deliveryTime,
		Status:          orderStatus.Name(),
		Items:           items,
		MerchantId:      merchantId,
		Amount:          amount,
	}, nil
}

func validateOrder(paymentMethodString, deliveryModeString, paymentModeString, cardToken string, payback float64) (enums.PaymentMethod, enums.DeliveryMode, enums.PaymentMode, error) {
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
	if payback < 0 {
		compositeError.Append(errors.PaybackShouldBePositiveError)
	}
	return paymentMethod, deliveryMode, paymentMode, compositeError.AsError()
}

func (order *Order) AddProduct(product *Product, quantity int, observation string) error {
	if product == nil {
		return errors.ProductNotFoundError
	}
	if product.IsUnAvailable() {
		return errors.ProductUnAvailableError
	}
	order.Items = append(order.Items, OrderItem{
		Id:          value_object.NewID(),
		ProductId:   product.Id,
		Quantity:    quantity,
		Observation: observation,
		Price:       product.Price,
		Status:      OrderItemStatusOk,
	})
	return nil
}

func (order *Order) Approve() error {
	err := order.state().Approve()
	if err != nil {
		return err
	}
	order.Status = status.ApprovedStatus.Name()
	return nil
}

func (order *Order) MarkAsInProgress() error {
	err := order.state().MarkAsInProgress()
	if err != nil {
		return err
	}
	order.Status = status.InProgressStatus.Name()
	return nil
}

func (order *Order) MarkAsAwaitingDelivery() error {
	if order.DeliveryMode.IsWithdraw() {
		return errors.OrderDeliveryModeIsWithdrawError
	}
	err := order.state().MarkAsInAwaitingDelivery()
	if err != nil {
		return err
	}
	order.Status = status.AwaitingDeliveryStatus.Name()
	return nil
}

func (order *Order) MarkAsAwaitingWithdraw() error {
	if order.DeliveryMode.IsDelivery() {
		return errors.OrderDeliveryModeIsDeliveryError
	}
	err := order.state().MarkAsInAwaitingWithdraw()
	if err != nil {
		return err
	}
	order.Status = status.AwaitingWithdrawStatus.Name()
	return nil
}

func (order *Order) MarkAsDelivering() error {
	err := order.state().MarkAsInDelivering()
	if err != nil {
		return err
	}
	order.Status = status.DeliveringStatus.Name()
	return nil
}

func (order *Order) Finish() error {
	err := order.state().Finish()
	if err != nil {
		return err
	}
	order.Status = status.FinishedStatus.Name()
	return nil
}

func (order *Order) Reject() error {
	err := order.state().Reject()
	if err != nil {
		return err
	}
	order.Status = status.RejectedStatus.Name()
	return nil
}

func (order *Order) Cancel() error {
	err := order.state().Cancel()
	if err != nil {
		return err
	}
	order.Status = status.CancelledStatus.Name()
	return nil
}

func (order *Order) GetStatus() string {
	return order.Status
}

func (order *Order) IsPixInApp() bool {
	return order.PaymentMethod.IsPix() && order.PaymentMode.IsInApp()
}

func (order *Order) IsAwaitingPayment() bool {
	return order.state() == status.AwaitingPaymentStatus
}

func (order *Order) IsAwaitingApproval() bool {
	return order.state() == status.AwaitingApprovalStatus
}

func (order *Order) AwaitApproval() error {
	if order.state() != status.AwaitingPaymentStatus {
		return status.OperationNotAllowedError
	}
	order.Status = status.AwaitingApprovalStatus.Name()
	return nil
}

func (order *Order) PaymentFailed() error {
	if order.state() != status.AwaitingPaymentStatus {
		return status.OperationNotAllowedError
	}
	order.Status = status.PaymentFailedStatus.Name()
	return nil
}

func (order *Order) state() status.Status {
	for _, state := range status.AllStatus {
		if state.Name() == order.Status {
			return state
		}
	}
	return nil
}
