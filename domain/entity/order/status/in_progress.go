package status

import (
	"time"
)

type InProgress struct {
}

func NewInProgress() Status {
	return &InProgress{}
}

func (status *InProgress) Name() string {
	return "InProgress"
}

func (status *InProgress) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *InProgress) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *InProgress) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*History, error) {
	return &History{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *InProgress) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *InProgress) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *InProgress) MarkAsInAwaitingWithdraw(at time.Time, by string) (*History, error) {
	return &History{
		Time:      at,
		Status:    "Delivering",
		ChangedBy: by,
	}, nil
}

func (status *InProgress) MarkAsInAwaitingDelivery(at time.Time, by string) (*History, error) {
	return &History{
		Time:      at,
		Status:    "AwaitingDelivery",
		ChangedBy: by,
	}, nil
}

func (status *InProgress) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
