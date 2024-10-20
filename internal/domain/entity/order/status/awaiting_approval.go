package status

type AwaitingApproval struct {
}

func NewAwaitingApproval() Status {
	return &AwaitingApproval{}
}

func (status *AwaitingApproval) Name() string {
	return "AwaitingApproval"
}

func (status *AwaitingApproval) Approve() error {
	return nil
}

func (status *AwaitingApproval) Reject() error {
	return nil
}

func (status *AwaitingApproval) Cancel() error {
	return OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInProgress() error {
	return OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInDelivering() error {
	return OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInAwaitingWithdraw() error {
	return OperationNotAllowedError
}

func (status *AwaitingApproval) MarkAsInAwaitingDelivery() error {
	return OperationNotAllowedError
}

func (status *AwaitingApproval) Finish() error {
	return OperationNotAllowedError
}
