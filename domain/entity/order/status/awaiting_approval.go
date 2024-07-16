package status

import (
	"bom-pedido-api/domain/entity/order"
	"time"
)

type AwaitingApproval struct {
}

func NewAwaitingApproval() order.Status {
	return &AwaitingApproval{}
}

func (status *AwaitingApproval) Name() string {
	return "AwaitingApproval"
}

func (status *AwaitingApproval) Approve(approvedAt time.Time, approvedBy string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      approvedAt,
		Status:    "Approved",
		ChangedBy: approvedBy,
	}, nil
}

func (status *AwaitingApproval) Reject(at time.Time, by string, reason string) (*order.StatusHistory, error) {
	return &order.StatusHistory{
		Time:      at,
		Status:    "Rejected",
		ChangedBy: by,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *AwaitingApproval) Cancel(time.Time, string, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInProgress(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInDelivering(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInAwaitingWithdraw(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInAwaitingDelivery(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}

func (status *AwaitingApproval) Finish(time.Time, string) (*order.StatusHistory, error) {
	return nil, order.OperationNotAllowedError
}
