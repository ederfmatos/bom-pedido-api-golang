package status

import (
	"time"
)

type AwaitingDelivery struct {
}

func NewAwaitingDelivery() Status {
	return &AwaitingDelivery{}
}

func (status *AwaitingDelivery) Name() string {
	return "AwaitingDelivery"
}

func (status *AwaitingDelivery) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingDelivery) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingDelivery) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*History, error) {
	return &History{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      reason,
	}, nil
}

func (status *AwaitingDelivery) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingDelivery) MarkAsInDelivering(at time.Time, by string) (*History, error) {
	return &History{
		Time:      at,
		Status:    "Delivering",
		ChangedBy: by,
	}, nil
}

func (status *AwaitingDelivery) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingDelivery) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingDelivery) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
