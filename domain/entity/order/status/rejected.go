package status

import (
	"time"
)

type Rejected struct {
}

func NewRejected() Status {
	return &Rejected{}
}

func (status *Rejected) Name() string {
	return "Rejected"
}

func (status *Rejected) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Rejected) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Rejected) Cancel(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Rejected) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Rejected) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Rejected) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Rejected) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Rejected) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
