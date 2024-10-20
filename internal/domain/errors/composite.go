package errors

import "strings"

type CompositeError struct {
	Errors []error
}

func NewCompositeError() *CompositeError {
	return &CompositeError{Errors: []error{}}
}

func NewCompositeWithError(err ...error) *CompositeError {
	return &CompositeError{Errors: err}
}

func (composite *CompositeError) Append(err error) {
	if err != nil {
		composite.Errors = append(composite.Errors, err)
	}
}

func (composite *CompositeError) AppendError(err error) {
	if err != nil {
		composite.Errors = append(composite.Errors, err)
	}
}

func (composite *CompositeError) AsError() error {
	if composite.HasError() {
		return composite
	}
	return nil
}

func (composite *CompositeError) HasError() bool {
	return len(composite.Errors) > 0
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
