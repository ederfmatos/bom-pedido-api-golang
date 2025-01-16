package google_authenticate_customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/application/token"
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type (
	UseCase struct {
		googleGateway      gateway.GoogleGateway
		customerRepository repository.CustomerRepository
		tokenManager       token.Manager
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
		googleGateway:      factory.GoogleGateway,
		customerRepository: factory.CustomerRepository,
		tokenManager:       factory.TokenManager,
	}
}

func (useCase *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
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
	return &Output{Token: customerToken}, nil
}
