package status

type Rejected struct {
}

func NewRejected() Status {
	return &Rejected{}
}

func (status *Rejected) Name() string {
	return "Rejected"
}

func (status *Rejected) Approve() error {
	return OperationNotAllowedError
}

func (status *Rejected) Reject() error {
	return OperationNotAllowedError
}

func (status *Rejected) Cancel() error {
	return OperationNotAllowedError
}

func (status *Rejected) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *Rejected) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *Rejected) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *Rejected) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *Rejected) Finish() error {
	return OperationNotAllowedError
}
