package errors

type DomainError struct {
	message string
}

func New(message string) *DomainError {
	return &DomainError{message: message}
}

func (error *DomainError) Error() string {
	return error.message
}
