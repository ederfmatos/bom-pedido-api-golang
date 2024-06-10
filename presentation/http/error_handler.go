package http

import (
	domainError "bom-pedido-api/domain/errors"
	"errors"
	"github.com/labstack/echo/v4"
)

func HandleError(err error, c echo.Context) {
	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	var errorResponse ErrorResponse
	var compositeError *domainError.CompositeError
	if errors.As(err, &compositeError) {
		errorResponse = ErrorResponse{Errors: compositeError.GetErrors()}
	} else {
		errorResponse = ErrorResponse{Errors: []string{err.Error()}}
	}
	_ = c.JSON(400, errorResponse)
}
