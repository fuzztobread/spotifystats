package handlers

import (

	"github.com/labstack/echo/v4"
)

func ServeDashboard(c echo.Context) error {
	return c.File("web/static/index.html")
}
