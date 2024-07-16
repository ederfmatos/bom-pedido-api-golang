package status

import (
	"time"
)

type Cancelled struct {
}

func NewCancelled() Status {
	return &Cancelled{}
}

func (status *Cancelled) Name() string {
	return "Cancelled"
}

func (status *Cancelled) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Cancelled) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Cancelled) Cancel(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Cancelled) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Cancelled) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Cancelled) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Cancelled) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Cancelled) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
