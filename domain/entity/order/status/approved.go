package status

import (
	"time"
)

type Approved struct {
}

func NewApproved() Status {
	return &Approved{}
}

func (status *Approved) Name() string {
	return "Approved"
}

func (status *Approved) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Approved) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Approved) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*History, error) {
	return &History{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      reason,
	}, nil
}

func (status *Approved) MarkAsInProgress(at time.Time, by string) (*History, error) {
	return &History{
		Time:      at,
		Status:    "InProgress",
		ChangedBy: by,
	}, nil
}

func (status *Approved) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Approved) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Approved) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Approved) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
