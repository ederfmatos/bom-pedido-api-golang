package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type InProgress struct {
}

func NewInProgress() order.Status {
	return &InProgress{}
}

func (status *InProgress) Name() string {
	return "InProgress"
}

func (status *InProgress) Approve(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *InProgress) Reject(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *InProgress) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *InProgress) MarkAsInProgress(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *InProgress) MarkAsInDelivering(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *InProgress) MarkAsInAwaitingWithdraw(at time.Time, by string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      at,
		Status:    "Delivering",
		ChangedBy: by,
	}, nil
}

func (status *InProgress) MarkAsInAwaitingDelivery(at time.Time, by string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      at,
		Status:    "AwaitingDelivery",
		ChangedBy: by,
	}, nil
}

func (status *InProgress) Finish(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}
