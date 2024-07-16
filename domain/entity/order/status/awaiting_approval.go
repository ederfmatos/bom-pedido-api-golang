package status

import (
	"time"
)

type AwaitingApproval struct {
}

func NewAwaitingApproval() Status {
	return &AwaitingApproval{}
}

func (status *AwaitingApproval) Name() string {
	return "AwaitingApproval"
}

func (status *AwaitingApproval) Approve(approvedAt time.Time, approvedBy string) (*History, error) {
	return &History{
		Time:      approvedAt,
		Status:    "Approved",
		ChangedBy: approvedBy,
	}, nil
}

func (status *AwaitingApproval) Reject(at time.Time, by string, reason string) (*History, error) {
	return &History{
		Time:      at,
		Status:    "Rejected",
		ChangedBy: by,
		Data:      map[string]string{"reason": reason},
	}, nil
}

func (status *AwaitingApproval) Cancel(time.Time, string, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInProgress(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInDelivering(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInAwaitingWithdraw(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInAwaitingDelivery(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}

func (status *AwaitingApproval) Finish(time.Time, string) (*History, error) {
	return nil, OperationNotAllowedError
}
