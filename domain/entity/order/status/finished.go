package status

type Finished struct {
}

func NewFinished() Status {
	return &Finished{}
}

func (status *Finished) Name() string {
	return "Finished"
}

func (status *Finished) Approve() error {
	return OperationNotAllowedError
}

func (status *Finished) Reject() error {
	return OperationNotAllowedError
}

func (status *Finished) Cancel() error {
	return OperationNotAllowedError
}

func (status *Finished) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *Finished) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *Finished) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *Finished) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *Finished) Finish() error {
	return OperationNotAllowedError
}
