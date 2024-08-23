package get_customer

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"context"
)

type (
	UseCase struct {
		customerRepository repository.CustomerRepository
	}
	Input struct {
		Id string
	}
	Output struct {
		Name        string  `json:"name,omitempty"`
		Email       string  `json:"email,omitempty"`
		PhoneNumber *string `json:"phoneNumber,omitempty"`
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		customerRepository: factory.CustomerRepository,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	customer, err := useCase.customerRepository.FindById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.CustomerNotFoundError
	}
	return &Output{
		Name:        customer.Name,
		Email:       customer.GetEmail(),
		PhoneNumber: customer.GetPhoneNumber(),
	}, nil
}
