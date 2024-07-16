package order

import (
	"bom-pedido-api/domain/entity/order/status"
	"bom-pedido-api/domain/errors"
	"time"
)

var (
	OperationNotAllowedError = errors.New("operation not allowed")
	InvalidStatusError       = errors.New("invalid status")

	AwaitingApproval = status.NewAwaitingApproval()
	Approved         = status.NewApproved()
	InProgress       = status.NewInProgress()
	Rejected         = status.NewRejected()
	Cancelled        = status.NewCancelled()
	Delivering       = status.NewDelivering()
	AwaitingWithdraw = status.NewAwaitingWithdraw()
	AwaitingDelivery = status.NewAwaitingDelivery()
	Finished         = status.NewFinished()
)

type (
	Status interface {
		Name() string
		Approve(approvedAt time.Time, approvedBy string) (*StatusHistory, error)
		Reject(rejectedAt time.Time, rejectedBy string, reason string) (*StatusHistory, error)
		Cancel(cancelledAt time.Time, cancelledBy string, reason string) (*StatusHistory, error)
		MarkAsInProgress(at time.Time, by string) (*StatusHistory, error)
		MarkAsInDelivering(at time.Time, by string) (*StatusHistory, error)
		MarkAsInAwaitingWithdraw(at time.Time, by string) (*StatusHistory, error)
		MarkAsInAwaitingDelivery(at time.Time, by string) (*StatusHistory, error)
		Finish(at time.Time, by string) (*StatusHistory, error)
	}
	StatusHistory struct {
		Time      time.Time
		Status    string
		ChangedBy string
		Data      map[string]string
	}
)

func ParseStatus(value string) (Status, error) {
	switch value {
	case AwaitingApproval.Name():
		return AwaitingApproval, nil
	case Approved.Name():
		return Approved, nil
	case InProgress.Name():
		return InProgress, nil
	case Rejected.Name():
		return Rejected, nil
	case Cancelled.Name():
		return Cancelled, nil
	case Delivering.Name():
		return Delivering, nil
	case AwaitingWithdraw.Name():
		return AwaitingWithdraw, nil
	case AwaitingDelivery.Name():
		return AwaitingDelivery, nil
	case Finished.Name():
		return Finished, nil
	default:
		return nil, InvalidStatusError
	}
}
