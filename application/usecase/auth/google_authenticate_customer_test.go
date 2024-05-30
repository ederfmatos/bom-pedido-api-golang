package auth

import (
	"bom-pedido-api/application/factory"
	gateway2 "bom-pedido-api/application/gateway"
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/repository"
	"bom-pedido-api/infra/token"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ShouldReturnsErrorIfGoogleReturnsError(t *testing.T) {
	googleGateway := gateway.NewFakeGoogleGateway()
	customerRepository := repository.NewCustomerMemoryRepository()
	manager := token.NewFakeCustomerTokenManager()

	applicationFactory := &factory.ApplicationFactory{
		GatewayFactory:    &factory.GatewayFactory{GoogleGateway: googleGateway},
		RepositoryFactory: &factory.RepositoryFactory{CustomerRepository: customerRepository},
		TokenFactory:      &factory.TokenFactory{CustomerTokenManager: manager},
	}

	googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(applicationFactory)
	input := Input{Token: "error", Context: context.Background()}
	googleGateway.On("GetUserByToken", "error").Return(nil, errors.New("any token"))

	output, err := googleAuthenticateCustomerUseCase.Execute(input)

	assert.Error(t, err, "expected error when using GoogleAuthenticateCustomerUseCase")
	assert.Nil(t, output)
}

func Test_ShouldReturnsErrorIfCustomerIsInvalid(t *testing.T) {
	googleGateway := gateway.NewFakeGoogleGateway()

	customerRepository := repository.NewCustomerMemoryRepository()
	manager := token.NewFakeCustomerTokenManager()

	applicationFactory := &factory.ApplicationFactory{
		GatewayFactory:    &factory.GatewayFactory{GoogleGateway: googleGateway},
		RepositoryFactory: &factory.RepositoryFactory{CustomerRepository: customerRepository},
		TokenFactory:      &factory.TokenFactory{CustomerTokenManager: manager},
	}
	googleUser := &gateway2.GoogleUser{
		Name:  faker.Name(),
		Email: "invalid_email",
	}

	googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(applicationFactory)
	input := Input{Token: "token", Context: context.Background()}
	googleGateway.On("GetUserByToken", "token").Return(googleUser, nil).Once()

	output, err := googleAuthenticateCustomerUseCase.Execute(input)

	assert.Error(t, err, "expected error when using GoogleAuthenticateCustomerUseCase")
	assert.ErrorIs(t, err, errors.InvalidEmailError)
	assert.Nil(t, output)
}

func Test_Should_Create_Customer_If_Does_Not_Exists(t *testing.T) {
	googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(factory.NewTestApplicationFactory())
	input := Input{Token: "google Token", Context: context.Background()}

	output, err := googleAuthenticateCustomerUseCase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "token", output.Token)
}

func Test_Should_Authenticate_Existent_Customer(t *testing.T) {
	testFactory := factory.NewTestApplicationFactory()
	googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(testFactory)
	customer, _ := entity.NewCustomer(faker.Name(), faker.Email())
	_ = testFactory.CustomerRepository.Create(context.TODO(), customer)
	input := Input{Token: "google Token", Context: context.Background()}

	output, err := googleAuthenticateCustomerUseCase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "token", output.Token)
}
