package order

import (
	"bom-pedido-api/domain/entity/order/status"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"time"
)

type ItemStatus string

const (
	ItemStatusOk ItemStatus = "OK"
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
		Payback         float64
		Code            int32
		DeliveryTime    time.Time
		state           status.Status
		Items           []Item
		MerchantId      string
		Amount          float64
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

func New(customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken string, payback, amount float64, deliveryTime time.Time, merchantId string) (*Order, error) {
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
		CreatedAt:       time.Now(),
		CreditCardToken: creditCardToken,
		Payback:         payback,
		Code:            0,
		DeliveryTime:    deliveryTime,
		state:           state,
		Items:           make([]Item, 0),
		MerchantId:      merchantId,
		Amount:          amount,
	}, nil
}

func Restore(Id, customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken, orderStatusString string, createdAt time.Time, payback, amount float64, code int32, deliveryTime time.Time, items []Item, merchantId string) (*Order, error) {
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
		state:           orderStatus,
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

func (order *Order) AddProduct(product *product.Product, quantity int, observation string) error {
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
		Status:      ItemStatusOk,
	})
	return nil
}

func (order *Order) Approve() error {
	err := order.state.Approve()
	if err != nil {
		return err
	}
	order.state = status.ApprovedStatus
	return nil
}

func (order *Order) MarkAsInProgress() error {
	err := order.state.MarkAsInProgress()
	if err != nil {
		return err
	}
	order.state = status.InProgressStatus
	return nil
}

func (order *Order) MarkAsAwaitingDelivery() error {
	if order.DeliveryMode.IsWithdraw() {
		return errors.OrderDeliveryModeIsWithdrawError
	}
	err := order.state.MarkAsInAwaitingDelivery()
	if err != nil {
		return err
	}
	order.state = status.AwaitingDeliveryStatus
	return nil
}

func (order *Order) MarkAsAwaitingWithdraw() error {
	if order.DeliveryMode.IsDelivery() {
		return errors.OrderDeliveryModeIsDeliveryError
	}
	err := order.state.MarkAsInAwaitingWithdraw()
	if err != nil {
		return err
	}
	order.state = status.AwaitingWithdrawStatus
	return nil
}

func (order *Order) MarkAsDelivering() error {
	err := order.state.MarkAsInDelivering()
	if err != nil {
		return err
	}
	order.state = status.DeliveringStatus
	return nil
}

func (order *Order) Finish() error {
	err := order.state.Finish()
	if err != nil {
		return err
	}
	order.state = status.FinishedStatus
	return nil
}

func (order *Order) Reject() error {
	err := order.state.Reject()
	if err != nil {
		return err
	}
	order.state = status.RejectedStatus
	return nil
}

func (order *Order) Cancel() error {
	err := order.state.Cancel()
	if err != nil {
		return err
	}
	order.state = status.CancelledStatus
	return nil
}

func (order *Order) GetStatus() string {
	return order.state.Name()
}

func (order *Order) IsPixInApp() bool {
	return order.PaymentMethod.IsPix() && order.PaymentMode.IsInApp()
}

func (order *Order) IsAwaitingPayment() bool {
	return order.state == status.AwaitingPaymentStatus
}

func (order *Order) IsAwaitingApproval() bool {
	return order.state == status.AwaitingApprovalStatus
}

func (order *Order) AwaitApproval() error {
	if order.state != status.AwaitingPaymentStatus {
		return status.OperationNotAllowedError
	}
	order.state = status.AwaitingApprovalStatus
	return nil
}

func (order *Order) PaymentFailed() error {
	if order.state != status.AwaitingPaymentStatus {
		return status.OperationNotAllowedError
	}
	order.state = status.PaymentFailedStatus
	return nil
}
