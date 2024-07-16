package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type Rejected struct {
}

func NewRejected() order.Status {
	return &Rejected{}
}

func (status *Rejected) Name() string {
	return "Rejected"
}

func (status *Rejected) Approve(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Rejected) Reject(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Rejected) Cancel(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Rejected) MarkAsInProgress(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Rejected) MarkAsInDelivering(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Rejected) MarkAsInAwaitingWithdraw(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Rejected) MarkAsInAwaitingDelivery(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *Rejected) Finish(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}
