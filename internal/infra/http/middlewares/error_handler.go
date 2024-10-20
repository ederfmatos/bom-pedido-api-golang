package middlewares

import (
	domainError "bom-pedido-api/internal/domain/errors"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func HandleError(err error, c echo.Context) {
	if c.Response().Header().Get("X-Error") != "" {
		return
	}
	slog.Error("Ocorreu um erro na requisição", slog.String("request", fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL)), slog.String("error", err.Error()))
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
