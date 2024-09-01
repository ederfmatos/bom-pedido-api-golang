package status

type Cancelled struct {
}

func NewCancelled() Status {
	return &Cancelled{}
}

func (status *Cancelled) Name() string {
	return "Cancelled"
}

func (status *Cancelled) Approve() error {
	return OperationNotAllowedError
}

func (status *Cancelled) Reject() error {
	return OperationNotAllowedError
}

func (status *Cancelled) Cancel() error {
	return OperationNotAllowedError
}

func (status *Cancelled) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *Cancelled) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *Cancelled) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *Cancelled) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *Cancelled) Finish() error {
	return OperationNotAllowedError
}
