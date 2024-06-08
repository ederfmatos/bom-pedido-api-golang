package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Health(t *testing.T) {
	instance := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	context := instance.NewContext(request, response)

	err := HandleHealth(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)

	var output HealthOutput
	_ = json.NewDecoder(response.Body).Decode(&output)
	assert.Equal(t, true, output.Ok)
}
