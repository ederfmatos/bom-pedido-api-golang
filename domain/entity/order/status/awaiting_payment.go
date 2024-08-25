package status

import (
	"time"
)

type AwaitingPayment struct {
}

func NewAwaitingPayment() Status {
	return &AwaitingPayment{}
}

func (status *AwaitingPayment) Name() string {
	return "AwaitingPayment"
}

func (status *AwaitingPayment) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingPayment) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingPayment) Cancel(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingPayment) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
