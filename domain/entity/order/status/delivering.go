package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type Delivering struct {
}

func NewDelivering() order.Status {
	return &Delivering{}
}

func (status *Delivering) Name() string {
	return "Delivering"
}

func (status *Delivering) Approve(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Delivering) Reject(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Delivering) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *Delivering) MarkAsInProgress(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Delivering) MarkAsInDelivering(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Delivering) MarkAsInAwaitingWithdraw(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Delivering) MarkAsInAwaitingDelivery(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Delivering) Finish(at time.Time, by string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      at,
		Status:    "Finished",
		ChangedBy: by,
	}, nil
}
