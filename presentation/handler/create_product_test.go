package handler

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateProduct(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()

	e := echo.New()
	body := createProductRequest{
		Name:        faker.Name(),
		Description: faker.Word(),
		Price:       10.0,
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		panic(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/v1/products", &buffer)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	echoContext := e.NewContext(request, response)

	err = HandleCreateProduct(applicationFactory)(echoContext)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)

	var output usecase.CreateProductOutput
	_ = json.NewDecoder(response.Body).Decode(&output)
	assert.NotEmpty(t, output.ID)

	savedProduct, err := applicationFactory.ProductRepository.FindById(context.Background(), output.ID)
	assert.NoError(t, err)
	assert.NotNil(t, savedProduct)
	assert.Equal(t, body.Name, savedProduct.Name)
	assert.Equal(t, body.Description, savedProduct.Description)
	assert.Equal(t, body.Price, savedProduct.Price)
	assert.True(t, savedProduct.IsActive())
}
