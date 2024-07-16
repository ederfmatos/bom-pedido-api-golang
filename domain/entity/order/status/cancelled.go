package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type Cancelled struct {
}

func NewCancelled() order.Status {
	return &Cancelled{}
}

func (status *Cancelled) Name() string {
	return "Cancelled"
}

func (status *Cancelled) Approve(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Cancelled) Reject(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Cancelled) Cancel(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Cancelled) MarkAsInProgress(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Cancelled) MarkAsInDelivering(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Cancelled) MarkAsInAwaitingWithdraw(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Cancelled) MarkAsInAwaitingDelivery(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Cancelled) Finish(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}
