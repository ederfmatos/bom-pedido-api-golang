package status

import (
	"time"
)

type Delivering struct {
}

func NewDelivering() Status {
	return &Delivering{}
}

func (status *Delivering) Name() string {
	return "Delivering"
}

func (status *Delivering) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Delivering) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Delivering) Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*History, error) {
	return &History{
		Time:      cancelledAt,
		Status:    "Cancelled",
		ChangedBy: cancelledBy,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *Delivering) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Delivering) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Delivering) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Delivering) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Delivering) Finish(at time.Time, by string) (*History, error) {
	return &History{
		Time:      at,
		Status:    "Finished",
		ChangedBy: by,
	}, nil
}
