package errors

import "strings"

type CompositeError struct {
	Errors []error
}

func NewCompositeError() *CompositeError {
	return &CompositeError{Errors: []error{}}
}

func NewCompositeWithError(err error) *CompositeError {
	return &CompositeError{Errors: []error{err}}
}

func (composite *CompositeError) Append(err *DomainError) {
	composite.Errors = append(composite.Errors, err)
}

func (composite *CompositeError) AsError() error {
	if len(composite.Errors) == 0 {
		return nil
	}
	return composite
}

func (composite *CompositeError) Error() string {
	errs := make([]string, len(composite.Errors))
	for i, err := range composite.Errors {
		errs[i] = err.Error()
	}
	return strings.Join(errs, "\n")
}

func (composite *CompositeError) GetErrors() []string {
	errs := make([]string, len(composite.Errors))
	for i, err := range composite.Errors {
		errs[i] = err.Error()
	}
	return errs
}
