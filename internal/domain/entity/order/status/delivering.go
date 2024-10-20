package status

type Delivering struct {
}

func NewDelivering() Status {
	return &Delivering{}
}

func (status *Delivering) Name() string {
	return "Delivering"
}

func (status *Delivering) Approve() error {
	return OperationNotAllowedError
}

func (status *Delivering) Reject() error {
	return OperationNotAllowedError
}

func (status *Delivering) Cancel() error {
	return nil
}

func (status *Delivering) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *Delivering) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *Delivering) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *Delivering) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *Delivering) Finish() error {
	return nil
}
