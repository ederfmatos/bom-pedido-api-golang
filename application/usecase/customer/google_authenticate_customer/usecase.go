package google_authenticate_customer

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/application/token"
	customerEntity "bom-pedido-api/domain/entity/customer"
	"context"
)

type (
	UseCase struct {
		googleGateway        gateway.GoogleGateway
		customerRepository   repository.CustomerRepository
		customerTokenManager token.CustomerTokenManager
	}
	Input struct {
		Token string
	}
	Output struct {
		Token string `json:"token"`
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		googleGateway:        factory.GoogleGateway,
		customerRepository:   factory.CustomerRepository,
		customerTokenManager: factory.CustomerTokenManager,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	googleUser, err := useCase.googleGateway.GetUserByToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}
	customer, err := useCase.customerRepository.FindByEmail(ctx, googleUser.Email)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		customer, err = customerEntity.New(googleUser.Name, googleUser.Email)
		if err != nil {
			return nil, err
		}
		err := useCase.customerRepository.Create(ctx, customer)
		if err != nil {
			return nil, err
		}
	}
	customerToken, err := useCase.customerTokenManager.Encrypt(ctx, customer.Id)
	return &Output{Token: customerToken}, err
}
