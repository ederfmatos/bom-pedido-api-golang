package order

import (
	"bom-pedido-api/domain/errors"
	"time"
)

var (
	AlreadyApprovedError  = errors.New("order already approved")
	AlreadyCancelledError = errors.New("order already cancelled")
	AlreadyRejectedError  = errors.New("order already rejected")
	AlreadyFinishedError  = errors.New("order already finished")
	ApprovalNotAllowed    = errors.New("approval not allowed")
	InvalidStatusError    = errors.New("invalid status")

	AwaitingApproval = newStatus("AWAITING_APPROVAL")
	Approved         = newStatus("APPROVED")
	InProgress       = newStatus("IN_PROGRESS")
	Rejected         = newStatus("REJECTED")
	Cancelled        = newStatus("CANCELLED")
	Delivering       = newStatus("DELIVERING")
	AwaitingWithdraw = newStatus("AWAITING_WITHDRAW")
	AwaitingDelivery = newStatus("AWAITING_DELIVERY")
	Finished         = newStatus("FINISHED")
)

type (
	Status struct {
		name string
	}
	StatusHistory struct {
		Time      time.Time
		Status    *Status
		ChangedBy string
	}
)

func newStatus(name string) *Status {
	return &Status{name}
}

func (s *Status) approve(approvedAt time.Time, approvedBy string) (*StatusHistory, error) {
	if s == AwaitingApproval {
		return &StatusHistory{
			Time:      approvedAt,
			Status:    Approved,
			ChangedBy: approvedBy,
		}, nil
	}
	switch s {
	case InProgress, Approved:
		return nil, AlreadyApprovedError
	case Cancelled:
		return nil, AlreadyCancelledError
	case Rejected:
		return nil, AlreadyRejectedError
	case Finished:
		return nil, AlreadyFinishedError
	default:
		return nil, ApprovalNotAllowed
	}
}

func ParseStatus(value string) (*Status, error) {
	switch value {
	case AwaitingApproval.name:
		return AwaitingApproval, nil
	case Approved.name:
		return Approved, nil
	case InProgress.name:
		return InProgress, nil
	case Rejected.name:
		return Rejected, nil
	case Cancelled.name:
		return Cancelled, nil
	case Delivering.name:
		return Delivering, nil
	case AwaitingWithdraw.name:
		return AwaitingWithdraw, nil
	case AwaitingDelivery.name:
		return AwaitingDelivery, nil
	case Finished.name:
		return Finished, nil
	default:
		return nil, InvalidStatusError
	}
}
