package status

type InProgress struct {
}

func NewInProgress() Status {
	return &InProgress{}
}

func (status *InProgress) Name() string {
	return "InProgress"
}

func (status *InProgress) Approve() error {
	return OperationNotAllowedError
}

func (status *InProgress) Reject() error {
	return OperationNotAllowedError
}

func (status *InProgress) Cancel() error {
	return nil
}

func (status *InProgress) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *InProgress) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *InProgress) MarkAsInAwaitingWithdraw() error {
	return nil
}

func (status *InProgress) MarkAsInAwaitingDelivery() error {
	return nil
}

func (status *InProgress) Finish() error {
	return OperationNotAllowedError
}
