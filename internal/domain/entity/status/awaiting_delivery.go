package status

type AwaitingDelivery struct {
}

func NewAwaitingDelivery() Status {
	return &AwaitingDelivery{}
}

func (status *AwaitingDelivery) Name() string {
	return "AwaitingDelivery"
}

func (status *AwaitingDelivery) Approve() error {
	return OperationNotAllowedError
}

func (status *AwaitingDelivery) Reject() error {
	return OperationNotAllowedError
}

func (status *AwaitingDelivery) Cancel() error {
	return nil
}

func (status *AwaitingDelivery) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *AwaitingDelivery) MarkAsInDelivering() error {
	return nil
}

func (status *AwaitingDelivery) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *AwaitingDelivery) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *AwaitingDelivery) Finish() error {
	return OperationNotAllowedError
}
