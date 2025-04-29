package api

import "github.com/labstack/echo/v4"

type Endpoint interface {
	Setup(e *echo.Echo)
	Execute(c echo.Context) error
}
