package google_authenticate_customer

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/gateway/google"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/internal/infra/token"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGoogleAuthenticateCustomerUseCase_Execute(t *testing.T) {
	googleGateway := google.NewFakeGoogleGateway()
	customerRepository := repository.NewCustomerMemoryRepository()
	tokenManager := token.NewFakeCustomerTokenManager()

	googleUser := &gateway.GoogleUser{
		Name:  faker.Name(),
		Email: faker.Email(),
	}

	applicationFactory := &factory.ApplicationFactory{
		GatewayFactory:    &factory.GatewayFactory{GoogleGateway: googleGateway},
		RepositoryFactory: &factory.RepositoryFactory{CustomerRepository: customerRepository},
		TokenFactory:      &factory.TokenFactory{TokenManager: tokenManager},
	}

	t.Run("ShouldReturnsErrorIfGoogleReturnsError", func(t *testing.T) {
		googleAuthenticateCustomerUseCase := New(applicationFactory)
		input := Input{Token: "error", TenantId: faker.Word()}
		googleGateway.On("GetUserByToken", "error").Return(nil, errors.New("any token"))

		output, err := googleAuthenticateCustomerUseCase.Execute(context.Background(), input)

		require.Error(t, err, "expected error when using UseCase")
		require.Nil(t, output)
	})

	t.Run("ShouldReturnsErrorIfCustomerIsInvalid", func(t *testing.T) {
		googleUser := &gateway.GoogleUser{
			Name:  faker.Name(),
			Email: "invalid_email",
		}

		googleAuthenticateCustomerUseCase := New(applicationFactory)
		input := Input{Token: "token", TenantId: faker.Word()}
		googleGateway.On("GetUserByToken", "token").Return(googleUser, nil).Once()

		output, err := googleAuthenticateCustomerUseCase.Execute(context.Background(), input)

		require.Error(t, err, "expected error when using UseCase")
		require.ErrorIs(t, err, value_object.InvalidEmailError)
		require.Nil(t, output)
	})

	t.Run("Should_Create_Customer_If_Does_Not_Exists", func(t *testing.T) {
		googleGateway.On("GetUserByToken", "google Token").Return(googleUser, nil).Once()
		tokenManager.On("Encrypt", mock.Anything).Return("token", nil).Once()

		googleAuthenticateCustomerUseCase := New(applicationFactory)
		input := Input{Token: "google Token", TenantId: faker.Word()}

		output, err := googleAuthenticateCustomerUseCase.Execute(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, output)
		require.Equal(t, "token", output.Token)
	})

	t.Run("Should_Authenticate_Existent_Customer", func(t *testing.T) {
		googleGateway.On("GetUserByToken", "google Token").Return(googleUser, nil).Once()
		tokenManager.On("Encrypt", mock.Anything).Return("token", nil).Once()

		googleAuthenticateCustomerUseCase := New(applicationFactory)
		customer, _ := entity.NewCustomer(faker.Name(), faker.Email(), faker.Word())
		_ = applicationFactory.CustomerRepository.Create(context.Background(), customer)
		input := Input{Token: "google Token", TenantId: faker.Word()}

		output, err := googleAuthenticateCustomerUseCase.Execute(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, output)
		require.Equal(t, "token", output.Token)
	})
}
