package health

import (
	"bom-pedido-api/internal/infra/json"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/testify/require"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Health(t *testing.T) {
	instance := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	container := test.NewContainer()
	redisClient := container.RedisClient

	err := Handle(redisClient, container.MongoDatabase())(instance.NewContext(request, response))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.Code)

	var output Output
	_ = json.Decode(request.Context(), response.Body, &output)
	require.Equal(t, true, output.Ok)
}
