package send_authentication_magic_link

import (
	"bom-pedido-api/domain/entity/admin"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/token"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	baseUrl := faker.URL()
	tokenManager := token.NewFakeCustomerTokenManager()
	eventEmitter := event.NewMockEventHandler()
	applicationFactory.TokenManager = tokenManager
	applicationFactory.EventEmitter = eventEmitter

	t.Run("it should return nil if admin does not exists", func(t *testing.T) {
		useCase := New(baseUrl, applicationFactory)
		input := Input{Email: faker.Email()}
		err := useCase.Execute(context.TODO(), input)
		assert.Nil(t, err)
	})

	t.Run("should return nil on success", func(t *testing.T) {
		anAdmin, _ := admin.New(faker.Name(), faker.Email(), faker.UUIDHyphenated())
		_ = applicationFactory.AdminRepository.Create(context.TODO(), anAdmin)

		magicLinkToken := faker.UUIDHyphenated()
		tokenManager.On("Encrypt", mock.Anything).Return(magicLinkToken, nil).Once()
		eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Once()

		useCase := New(baseUrl, applicationFactory)
		input := Input{Email: anAdmin.GetEmail()}

		ctx := context.TODO()
		err := useCase.Execute(ctx, input)
		assert.NoError(t, err)

		eventEmitter.AssertNumberOfCalls(t, "Emit", 1)
	})
}
