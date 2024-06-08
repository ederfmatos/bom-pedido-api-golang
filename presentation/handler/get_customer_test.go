package handler

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase"
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"encoding/json"
	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetCustomer(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("should success get customer", func(t *testing.T) {
		customer, _ := entity.NewCustomer(faker.Name(), faker.Email())
		_ = customer.SetPhoneNumber(faker.Phonenumber())
		_ = applicationFactory.CustomerRepository.Create(context.Background(), customer)

		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/customers/me", nil)
		response := httptest.NewRecorder()
		echoContext := e.NewContext(request, response)
		echoContext.Set("currentUserId", customer.Id)

		err := HandleGetAuthenticatedCustomer(applicationFactory)(echoContext)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)

		var output usecase.GetCustomerOutput
		_ = json.NewDecoder(response.Body).Decode(&output)
		assert.Equal(t, customer.Name, output.Name)
		assert.Equal(t, customer.GetEmail(), output.Email)
		assert.Equal(t, customer.GetPhoneNumber(), output.PhoneNumber)
	})

	t.Run("should return CustomerNotFoundError if customer does not exists", func(t *testing.T) {
		instance := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/customers/me", nil)
		response := httptest.NewRecorder()
		echoContext := instance.NewContext(request, response)
		echoContext.Set("currentUserId", value_object.NewID())

		err := HandleGetAuthenticatedCustomer(applicationFactory)(echoContext)
		assert.Equal(t, errors.CustomerNotFoundError, err)
	})
}
