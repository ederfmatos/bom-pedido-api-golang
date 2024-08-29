package create_product

import (
	"bom-pedido-api/application/usecase/product/create_product"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/json"
	"bom-pedido-api/infra/tenant"
	"bytes"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
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
	err := json.Encode(context.Background(), &buffer, body)
	if err != nil {
		panic(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/v1/products", &buffer)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	echoContext := e.NewContext(request, response)
	echoContext.Set(tenant.Id, faker.WORD)

	err = Handle(applicationFactory)(echoContext)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, response.Code)

	var output create_product.Output
	_ = json.Decode(request.Context(), response.Body, &output)
	require.NotEmpty(t, output.Id)

	savedProduct, err := applicationFactory.ProductRepository.FindById(context.TODO(), output.Id)
	require.NoError(t, err)
	require.NotNil(t, savedProduct)
	require.Equal(t, body.Name, savedProduct.Name)
	require.Equal(t, body.Description, savedProduct.Description)
	require.Equal(t, body.Price, savedProduct.Price)
	require.True(t, savedProduct.IsActive())
}
