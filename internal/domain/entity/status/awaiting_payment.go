package status

type AwaitingPayment struct {
}

func NewAwaitingPayment() Status {
	return &AwaitingPayment{}
}

func (status *AwaitingPayment) Name() string {
	return "AwaitingPayment"
}

func (status *AwaitingPayment) Approve() error {
	return OperationNotAllowedError
}

func (status *AwaitingPayment) Reject() error {
	return OperationNotAllowedError
}

func (status *AwaitingPayment) Cancel() error {
	return OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *AwaitingPayment) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *AwaitingPayment) Finish() error {
	return OperationNotAllowedError
}
