package middlewares

import (
	domainError "bom-pedido-api/domain/errors"
	"errors"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func HandleError(err error, c echo.Context) {
	if c.Response().Header().Get("X-Error") != "" {
		return
	}
	var errorResponse ErrorResponse
	var compositeError *domainError.CompositeError
	if errors.As(err, &compositeError) {
		errorResponse = ErrorResponse{Errors: compositeError.GetErrors()}
	} else {
		errorResponse = ErrorResponse{Errors: []string{err.Error()}}
	}
	c.Response().Header().Set("X-Error", err.Error())
	_ = c.JSON(400, errorResponse)
}
