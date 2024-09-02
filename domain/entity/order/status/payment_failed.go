package status

type PaymentFailed struct {
}

func NewPaymentFailed() Status {
	return &PaymentFailed{}
}

func (status *PaymentFailed) Name() string {
	return "PaymentFailed"
}

func (status *PaymentFailed) Approve() error {
	return OperationNotAllowedError
}

func (status *PaymentFailed) Reject() error {
	return OperationNotAllowedError
}

func (status *PaymentFailed) Cancel() error {
	return OperationNotAllowedError
}

func (status *PaymentFailed) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *PaymentFailed) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *PaymentFailed) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *PaymentFailed) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *PaymentFailed) Finish() error {
	return OperationNotAllowedError
}
