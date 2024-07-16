package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type AwaitingWithdraw struct {
}

func NewAwaitingWithdraw() order.Status {
	return &AwaitingWithdraw{}
}

func (status *AwaitingWithdraw) Name() string {
	return "AwaitingWithdraw"
}

func (status *AwaitingWithdraw) Approve(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingWithdraw) Reject(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingWithdraw) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *AwaitingWithdraw) MarkAsInProgress(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInDelivering(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInAwaitingWithdraw(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInAwaitingDelivery(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingWithdraw) Finish(at time.Time, by string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      at,
		Status:    "Finished",
		ChangedBy: by,
	}, nil
}
