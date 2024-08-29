package health

import (
	"bom-pedido-api/infra/json"
	"bom-pedido-api/infra/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Health(t *testing.T) {
	instance := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	container := test.NewContainer()
	database, mongoClient, redisClient := container.Database, container.MongoClient, container.RedisClient

	err := Handle(database, redisClient, mongoClient)(instance.NewContext(request, response))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.Code)

	var output Output
	_ = json.Decode(request.Context(), response.Body, &output)
	require.Equal(t, true, output.Ok)
}
