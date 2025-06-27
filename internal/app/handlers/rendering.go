package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func renderTemplate(c echo.Context, name string, model interface{}) error {
	return c.Render(http.StatusOK, name+".jet", model)
}

func errorHandler(c echo.Context, status int) error {
	return echo.NewHTTPError(status)
}
