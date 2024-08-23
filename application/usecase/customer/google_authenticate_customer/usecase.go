package google_authenticate_customer

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/application/token"
	"bom-pedido-api/domain/entity/customer"
	"context"
)

type (
	UseCase struct {
		googleGateway        gateway.GoogleGateway
		customerRepository   repository.CustomerRepository
		customerTokenManager token.CustomerTokenManager
	}
	Input struct {
		Token    string
		TenantId string
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
	aCustomer, err := useCase.customerRepository.FindByEmail(ctx, googleUser.Email, input.TenantId)
	if err != nil {
		return nil, err
	}
	if aCustomer == nil {
		aCustomer, err = customer.New(googleUser.Name, googleUser.Email, input.TenantId)
		if err != nil {
			return nil, err
		}
		err = useCase.customerRepository.Create(ctx, aCustomer)
		if err != nil {
			return nil, err
		}
	}
	customerToken, err := useCase.customerTokenManager.Encrypt(ctx, aCustomer.Id)
	return &Output{Token: customerToken}, err
}
