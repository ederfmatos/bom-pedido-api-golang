package status

type Approved struct {
}

func NewApproved() Status {
	return &Approved{}
}

func (status *Approved) Name() string {
	return "Approved"
}

func (status *Approved) Approve() error {
	return OperationNotAllowedError
}

func (status *Approved) Reject() error {
	return OperationNotAllowedError
}

func (status *Approved) Cancel() error {
	return nil
}

func (status *Approved) MarkAsInProgress() error {
	return nil
}

func (status *Approved) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *Approved) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *Approved) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *Approved) Finish() error {
	return OperationNotAllowedError
}
