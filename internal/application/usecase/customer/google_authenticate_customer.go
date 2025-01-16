package customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/application/token"
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type (
	GoogleAuthenticateCustomerUseCase struct {
		googleGateway      gateway.GoogleGateway
		customerRepository repository.CustomerRepository
		tokenManager       token.Manager
	}
	GoogleAuthenticateCustomerInput struct {
		Token    string
		TenantId string
	}
	GoogleAuthenticateCustomerOutput struct {
		Token string `json:"token"`
	}
)

func NewGoogleAuthenticateCustomer(factory *factory.ApplicationFactory) *GoogleAuthenticateCustomerUseCase {
	return &GoogleAuthenticateCustomerUseCase{
		googleGateway:      factory.GoogleGateway,
		customerRepository: factory.CustomerRepository,
		tokenManager:       factory.TokenManager,
	}
}

func (useCase *GoogleAuthenticateCustomerUseCase) Execute(ctx context.Context, input GoogleAuthenticateCustomerInput) (*GoogleAuthenticateCustomerOutput, error) {
	googleUser, err := useCase.googleGateway.GetUserByToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}
	customer, err := useCase.customerRepository.FindByEmail(ctx, googleUser.Email, input.TenantId)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		customer, err = entity.NewCustomer(googleUser.Name, googleUser.Email, input.TenantId)
		if err != nil {
			return nil, err
		}
		err = useCase.customerRepository.Create(ctx, customer)
		if err != nil {
			return nil, err
		}
	}
	tokenData := token.Data{
		Type:     "CUSTOMER",
		Id:       customer.Id,
		TenantId: input.TenantId,
	}
	customerToken, err := useCase.tokenManager.Encrypt(ctx, tokenData)
	if err != nil {
		return nil, err
	}
	return &GoogleAuthenticateCustomerOutput{Token: customerToken}, nil
}
