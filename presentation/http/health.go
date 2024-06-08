package http

import (
	"github.com/labstack/echo/v4"
)

type HealthOutput struct {
	Ok bool `json:"ok"`
}

func HandleHealth(context echo.Context) error {
	return context.JSON(200, HealthOutput{Ok: true})
}
