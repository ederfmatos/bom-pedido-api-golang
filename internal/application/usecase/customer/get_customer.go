package customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/errors"
	"context"
)

type (
	GetCustomerUseCase struct {
		customerRepository repository.CustomerRepository
	}
	GetCustomerInput struct {
		Id string
	}
	GetCustomerOutput struct {
		Name        string  `json:"name,omitempty"`
		Email       string  `json:"email,omitempty"`
		PhoneNumber *string `json:"phoneNumber,omitempty"`
	}
)

func NewGetCustomer(factory *factory.ApplicationFactory) *GetCustomerUseCase {
	return &GetCustomerUseCase{
		customerRepository: factory.CustomerRepository,
	}
}

func (useCase *GetCustomerUseCase) Execute(ctx context.Context, input GetCustomerInput) (*GetCustomerOutput, error) {
	customer, err := useCase.customerRepository.FindById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.CustomerNotFoundError
	}
	return &GetCustomerOutput{
		Name:        customer.Name,
		Email:       customer.GetEmail(),
		PhoneNumber: customer.GetPhoneNumber(),
	}, nil
}
