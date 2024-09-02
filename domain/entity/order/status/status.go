package status

import (
	"bom-pedido-api/domain/errors"
)

var (
	OperationNotAllowedError = errors.New("operation not allowed")
	InvalidStatusError       = errors.New("invalid status")

	AwaitingPaymentStatus  = NewAwaitingPayment()
	AwaitingApprovalStatus = NewAwaitingApproval()
	PaymentFailedStatus    = NewPaymentFailed()
	ApprovedStatus         = NewApproved()
	InProgressStatus       = NewInProgress()
	RejectedStatus         = NewRejected()
	CancelledStatus        = NewCancelled()
	DeliveringStatus       = NewDelivering()
	AwaitingWithdrawStatus = NewAwaitingWithdraw()
	AwaitingDeliveryStatus = NewAwaitingDelivery()
	FinishedStatus         = NewFinished()

	AllStatus = []Status{
		AwaitingPaymentStatus,
		AwaitingApprovalStatus,
		ApprovedStatus,
		InProgressStatus,
		RejectedStatus,
		CancelledStatus,
		DeliveringStatus,
		AwaitingWithdrawStatus,
		AwaitingDeliveryStatus,
		FinishedStatus,
		PaymentFailedStatus,
	}
)

type Status interface {
	Name() string
	Approve() error
	Reject() error
	Cancel() error
	MarkAsInProgress() error
	MarkAsInDelivering() error
	MarkAsInAwaitingWithdraw() error
	MarkAsInAwaitingDelivery() error
	Finish() error
}

func Parse(value string) (Status, error) {
	for _, status := range AllStatus {
		if status.Name() == value {
			return status, nil
		}
	}
	return nil, InvalidStatusError
}
