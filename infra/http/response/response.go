package response

import "github.com/labstack/echo/v4"

func Ok(context echo.Context, output interface{}, err error) error {
	if err != nil {
		return err
	}
	return context.JSON(200, output)
}

func Created(context echo.Context, output interface{}, err error) error {
	if err != nil {
		return err
	}
	return context.JSON(201, output)
}

func NoContent(context echo.Context, err error) error {
	if err != nil {
		return err
	}
	return context.JSON(204, nil)
}
