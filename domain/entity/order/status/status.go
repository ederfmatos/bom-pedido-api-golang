package status

import (
	"bom-pedido-api/domain/errors"
	"time"
)

var (
	OperationNotAllowedError = errors.New("operation not allowed")
	InvalidStatusError       = errors.New("invalid status")

	AwaitingApprovalStatus = NewAwaitingApproval()
	ApprovedStatus         = NewApproved()
	InProgressStatus       = NewInProgress()
	RejectedStatus         = NewRejected()
	CancelledStatus        = NewCancelled()
	DeliveringStatus       = NewDelivering()
	AwaitingWithdrawStatus = NewAwaitingWithdraw()
	AwaitingDeliveryStatus = NewAwaitingDelivery()
	FinishedStatus         = NewFinished()

	AllStatus = []Status{
		AwaitingApprovalStatus,
		ApprovedStatus,
		InProgressStatus,
		RejectedStatus,
		CancelledStatus,
		DeliveringStatus,
		AwaitingWithdrawStatus,
		AwaitingDeliveryStatus,
		FinishedStatus,
	}
)

type (
	Status interface {
		Name() string
		Approve(approvedAt time.Time, approvedBy string) (*History, error)
		Reject(rejectedAt time.Time, rejectedBy string, reason string) (*History, error)
		Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*History, error)
		MarkAsInProgress(at time.Time, by string) (*History, error)
		MarkAsInDelivering(at time.Time, by string) (*History, error)
		MarkAsInAwaitingWithdraw(at time.Time, by string) (*History, error)
		MarkAsInAwaitingDelivery(at time.Time, by string) (*History, error)
		Finish(at time.Time, by string) (*History, error)
	}
	History struct {
		Time      time.Time
		Status    string
		ChangedBy string
		Data      map[string]string
	}
)

func Parse(value string) (Status, error) {
	for _, status := range AllStatus {
		if status.Name() == value {
			return status, nil
		}
	}
	return nil, InvalidStatusError
}
