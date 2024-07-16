package status

import (
	"time"
)

type AwaitingWithdraw struct {
}

func NewAwaitingWithdraw() Status {
	return &AwaitingWithdraw{}
}

func (status *AwaitingWithdraw) Name() string {
	return "AwaitingWithdraw"
}

func (status *AwaitingWithdraw) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingWithdraw) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingWithdraw) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*History, error) {
	return &History{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *AwaitingWithdraw) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingWithdraw) Finish(at time.Time, by string) (*History, error) {
	return &History{
		Time:      at,
		Status:    "Finished",
		ChangedBy: by,
	}, nil
}
