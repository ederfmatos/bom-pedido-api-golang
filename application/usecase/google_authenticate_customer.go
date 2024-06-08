package usecase

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/application/token"
	"bom-pedido-api/domain/entity"
	"context"
)

type GoogleAuthenticateCustomerInput struct {
	Token   string
	Context context.Context
}

type GoogleAuthenticateCustomerOutput struct {
	Token string
}

type GoogleAuthenticateCustomerUseCase struct {
	googleGateway        gateway.GoogleGateway
	customerRepository   repository.CustomerRepository
	customerTokenManager token.CustomerTokenManager
}

func NewGoogleAuthenticateCustomerUseCase(factory *factory.ApplicationFactory) *GoogleAuthenticateCustomerUseCase {
	return &GoogleAuthenticateCustomerUseCase{
		googleGateway:        factory.GoogleGateway,
		customerRepository:   factory.CustomerRepository,
		customerTokenManager: factory.CustomerTokenManager,
	}
}

func (useCase GoogleAuthenticateCustomerUseCase) Execute(input GoogleAuthenticateCustomerInput) (*GoogleAuthenticateCustomerOutput, error) {
	googleUser, err := useCase.googleGateway.GetUserByToken(input.Token)
	if err != nil {
		return nil, err
	}
	customer, err := useCase.customerRepository.FindByEmail(input.Context, googleUser.Email)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		customer, err = entity.NewCustomer(googleUser.Name, googleUser.Email)
		if err != nil {
			return nil, err
		}
		err := useCase.customerRepository.Create(input.Context, customer)
		if err != nil {
			return nil, err
		}
	}
	customerToken, err := useCase.customerTokenManager.Encrypt(customer.Id)
	return &GoogleAuthenticateCustomerOutput{Token: customerToken}, err
}
