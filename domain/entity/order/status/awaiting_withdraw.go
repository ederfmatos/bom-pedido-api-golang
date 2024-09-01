package status

type AwaitingWithdraw struct {
}

func NewAwaitingWithdraw() Status {
	return &AwaitingWithdraw{}
}

func (status *AwaitingWithdraw) Name() string {
	return "AwaitingWithdraw"
}

func (status *AwaitingWithdraw) Approve() error {
	return OperationNotAllowedError
}

func (status *AwaitingWithdraw) Reject() error {
	return OperationNotAllowedError
}

func (status *AwaitingWithdraw) Cancel() error {
	return nil
}

func (status *AwaitingWithdraw) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *AwaitingWithdraw) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *AwaitingWithdraw) Finish() error {
	return nil
}
