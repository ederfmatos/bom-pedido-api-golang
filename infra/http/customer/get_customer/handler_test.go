package get_customer

import (
	"bom-pedido-api/application/usecase/customer/get_customer"
	"bom-pedido-api/domain/entity/customer"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/json"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetCustomer(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("should success get customer", func(t *testing.T) {
		customer, _ := customer.New(faker.Name(), faker.Email())
		_ = customer.SetPhoneNumber(faker.Phonenumber())
		_ = applicationFactory.CustomerRepository.Create(context.TODO(), customer)

		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/customers/me", nil)
		response := httptest.NewRecorder()
		echoContext := e.NewContext(request, response)
		echoContext.Set("customerId", customer.Id)

		err := Handle(applicationFactory)(echoContext)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)

		var output get_customer.Output
		_ = json.Decode(request.Context(), response.Body, &output)
		assert.Equal(t, customer.Name, output.Name)
		assert.Equal(t, customer.GetEmail(), output.Email)
		assert.Equal(t, customer.GetPhoneNumber(), output.PhoneNumber)
	})

	t.Run("should return CustomerNotFoundError if customer does not exists", func(t *testing.T) {
		instance := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/customers/me", nil)
		response := httptest.NewRecorder()
		echoContext := instance.NewContext(request, response)
		echoContext.Set("customerId", value_object.NewID())

		err := Handle(applicationFactory)(echoContext)
		assert.Equal(t, errors.CustomerNotFoundError, err)
	})
}
