package usecase

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
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGoogleAuthenticateCustomerUseCase_Execute(t *testing.T) {
	googleGateway := gateway.NewFakeGoogleGateway()
	customerRepository := repository.NewCustomerMemoryRepository()
	tokenManager := token.NewFakeCustomerTokenManager()

	googleUser := &gateway2.GoogleUser{
		Name:  faker.Name(),
		Email: faker.Email(),
	}

	applicationFactory := &factory.ApplicationFactory{
		GatewayFactory:    &factory.GatewayFactory{GoogleGateway: googleGateway},
		RepositoryFactory: &factory.RepositoryFactory{CustomerRepository: customerRepository},
		TokenFactory:      &factory.TokenFactory{CustomerTokenManager: tokenManager},
	}

	t.Run("ShouldReturnsErrorIfGoogleReturnsError", func(t *testing.T) {
		googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(applicationFactory)
		input := GoogleAuthenticateCustomerInput{Token: "error", Context: context.Background()}
		googleGateway.On("GetUserByToken", "error").Return(nil, errors.New("any token"))

		output, err := googleAuthenticateCustomerUseCase.Execute(input)

		assert.Error(t, err, "expected error when using GoogleAuthenticateCustomerUseCase")
		assert.Nil(t, output)
	})

	t.Run("ShouldReturnsErrorIfCustomerIsInvalid", func(t *testing.T) {
		googleUser := &gateway2.GoogleUser{
			Name:  faker.Name(),
			Email: "invalid_email",
		}

		googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(applicationFactory)
		input := GoogleAuthenticateCustomerInput{Token: "token", Context: context.Background()}
		googleGateway.On("GetUserByToken", "token").Return(googleUser, nil).Once()

		output, err := googleAuthenticateCustomerUseCase.Execute(input)

		assert.Error(t, err, "expected error when using GoogleAuthenticateCustomerUseCase")
		assert.ErrorIs(t, err, errors.InvalidEmailError)
		assert.Nil(t, output)
	})

	t.Run("Should_Create_Customer_If_Does_Not_Exists", func(t *testing.T) {

		googleGateway.On("GetUserByToken", "google Token").Return(googleUser, nil).Once()
		tokenManager.On("Encrypt", mock.Anything).Return("token", nil).Once()

		googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(applicationFactory)
		input := GoogleAuthenticateCustomerInput{Token: "google Token", Context: context.Background()}

		output, err := googleAuthenticateCustomerUseCase.Execute(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "token", output.Token)
	})

	t.Run("Should_Authenticate_Existent_Customer", func(t *testing.T) {
		googleGateway.On("GetUserByToken", "google Token").Return(googleUser, nil).Once()
		tokenManager.On("Encrypt", mock.Anything).Return("token", nil).Once()

		googleAuthenticateCustomerUseCase := NewGoogleAuthenticateCustomerUseCase(applicationFactory)
		customer, _ := entity.NewCustomer(faker.Name(), faker.Email())
		_ = applicationFactory.CustomerRepository.Create(context.TODO(), customer)
		input := GoogleAuthenticateCustomerInput{Token: "google Token", Context: context.Background()}

		output, err := googleAuthenticateCustomerUseCase.Execute(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "token", output.Token)
	})
}
