package usecase

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type GetCustomerInput struct {
	Id      string
	Context context.Context
}

type GetCustomerOutput struct {
	Name        string  `json:"name,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

type GetCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewGetCustomerUseCase(factory *factory.ApplicationFactory) *GetCustomerUseCase {
	return &GetCustomerUseCase{
		customerRepository: factory.CustomerRepository,
	}
}

func (useCase GetCustomerUseCase) Execute(input GetCustomerInput) (*GetCustomerOutput, error) {
	customer, err := useCase.customerRepository.FindById(input.Context, input.Id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, entity.CustomerNotFoundError
	}
	return &GetCustomerOutput{
		Name:        customer.Name,
		Email:       customer.GetEmail(),
		PhoneNumber: customer.GetPhoneNumber(),
	}, nil
}
