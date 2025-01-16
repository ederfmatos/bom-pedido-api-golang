package get_customer

import (
	"bom-pedido-api/internal/application/usecase/customer/get_customer"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/json"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetCustomer(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("should success get customer", func(t *testing.T) {
		customer, _ := entity.NewCustomer(faker.Name(), faker.Email(), faker.WORD)
		_ = customer.SetPhoneNumber(faker.Phonenumber())
		_ = applicationFactory.CustomerRepository.Create(context.Background(), customer)

		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/customers/me", nil)
		response := httptest.NewRecorder()
		echoContext := e.NewContext(request, response)
		echoContext.Set(middlewares.CustomerIdParam, customer.Id)

		err := Handle(applicationFactory)(echoContext)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.Code)

		var output get_customer.Output
		_ = json.Decode(request.Context(), response.Body, &output)
		require.Equal(t, customer.Name, output.Name)
		require.Equal(t, customer.GetEmail(), output.Email)
		require.Equal(t, customer.GetPhoneNumber(), output.PhoneNumber)
	})

	t.Run("should return CustomerNotFoundError if customer does not exists", func(t *testing.T) {
		instance := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/customers/me", nil)
		response := httptest.NewRecorder()
		echoContext := instance.NewContext(request, response)
		echoContext.Set(middlewares.CustomerIdParam, value_object.NewID())

		err := Handle(applicationFactory)(echoContext)
		require.Equal(t, errors.CustomerNotFoundError, err)
	})
}
