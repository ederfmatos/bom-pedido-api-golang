package http

import "github.com/labstack/echo/v4"

func HandleError(err error, c echo.Context) {
	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	errorResponse := ErrorResponse{
		Errors: []string{err.Error()},
	}
	_ = c.JSON(400, errorResponse)
}
