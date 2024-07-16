package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type Approved struct {
}

func NewApproved() order.Status {
	return &Approved{}
}

func (status *Approved) Name() string {
	return "Approved"
}

func (status *Approved) Approve(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Approved) Reject(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Approved) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *Approved) MarkAsInProgress(at time.Time, by string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      at,
		Status:    "InProgress",
		ChangedBy: by,
	}, nil
}

func (status *Approved) MarkAsInDelivering(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Approved) MarkAsInAwaitingWithdraw(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Approved) MarkAsInAwaitingDelivery(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Approved) Finish(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}
