package status

import (
	"time"
)

type Finished struct {
}

func NewFinished() Status {
	return &Finished{}
}

func (status *Finished) Name() string {
	return "Finished"
}

func (status *Finished) Approve(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Finished) Reject(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Finished) Cancel(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Finished) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Finished) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Finished) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Finished) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *Finished) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
