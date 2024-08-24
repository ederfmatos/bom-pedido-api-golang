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
		Payback         float64
		Code            int32
		DeliveryTime    time.Time
		state           status.Status
		Items           []Item
		History         []status.History
		TenantId        string
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

func New(customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken string, payback, amount float64, deliveryTime time.Time, tenantId string) (*Order, error) {
	paymentMethod, deliveryMode, paymentMode, err := validateOrder(paymentMethodString, deliveryModeString, paymentModeString, creditCardToken, payback)
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
		Payback:         payback,
		Code:            0,
		DeliveryTime:    deliveryTime,
		state:           status.AwaitingApprovalStatus,
		Items:           make([]Item, 0),
		History:         make([]status.History, 0),
		TenantId:        tenantId,
		Amount:          amount,
	}, nil
}

func Restore(
	Id, customerID, paymentMethodString, paymentModeString, deliveryModeString, creditCardToken, orderStatusString string,
	createdAt time.Time,
	payback, amount float64,
	code int32,
	deliveryTime time.Time,
	items []Item,
	history []status.History,
	tenantId string,
) (*Order, error) {
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
		History:         history,
		TenantId:        tenantId,
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
		Status:      OrderItemStatusOk,
	})
	return nil
}

func (order *Order) Approve(approvedAt time.Time, approvedBy string) error {
	history, err := order.state.Approve(approvedAt, approvedBy)
	if err != nil {
		return err
	}
	order.state = status.ApprovedStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) MarkAsInProgress(at time.Time, by string) error {
	history, err := order.state.MarkAsInProgress(at, by)
	if err != nil {
		return err
	}
	order.state = status.InProgressStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) MarkAsAwaitingDelivery(at time.Time, by string) error {
	if order.DeliveryMode.IsWithdraw() {
		return errors.OrderDeliveryModeIsWithdrawError
	}
	history, err := order.state.MarkAsInAwaitingDelivery(at, by)
	if err != nil {
		return err
	}
	order.state = status.AwaitingDeliveryStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) MarkAsAwaitingWithdraw(at time.Time, by string) error {
	if order.DeliveryMode.IsDelivery() {
		return errors.OrderDeliveryModeIsDeliveryError
	}
	history, err := order.state.MarkAsInAwaitingWithdraw(at, by)
	if err != nil {
		return err
	}
	order.state = status.AwaitingWithdrawStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) MarkAsDelivering(at time.Time, by string) error {
	history, err := order.state.MarkAsInDelivering(at, by)
	if err != nil {
		return err
	}
	order.state = status.DeliveringStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) Finish(at time.Time, by string) error {
	history, err := order.state.Finish(at, by)
	if err != nil {
		return err
	}
	order.state = status.FinishedStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) Reject(rejectedAt time.Time, rejectedBy, reason string) error {
	history, err := order.state.Reject(rejectedAt, rejectedBy, reason)
	if err != nil {
		return err
	}
	order.state = status.RejectedStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) Cancel(at time.Time, by, reason string) error {
	history, err := order.state.Cancel(at, by, reason)
	if err != nil {
		return err
	}
	order.state = status.CancelledStatus
	order.History = append(order.History, *history)
	return nil
}

func (order *Order) GetStatus() string {
	return order.state.Name()
}

func (order *Order) IsPixInApp() bool {
	return order.PaymentMethod == enums.PaymentMethodPix && order.PaymentMode == enums.PaymentModeInApp
}
