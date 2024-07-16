package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type AwaitingDelivery struct {
}

func NewAwaitingDelivery() order.Status {
	return &AwaitingDelivery{}
}

func (status *AwaitingDelivery) Name() string {
	return "AwaitingDelivery"
}

func (status *AwaitingDelivery) Approve(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingDelivery) Reject(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingDelivery) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *AwaitingDelivery) MarkAsInProgress(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingDelivery) MarkAsInDelivering(at time.Time, by string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      at,
		Status:    "Delivering",
		ChangedBy: by,
	}, nil
}

func (status *AwaitingDelivery) MarkAsInAwaitingWithdraw(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingDelivery) MarkAsInAwaitingDelivery(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingDelivery) Finish(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}
