package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "Safe and sound ❤"})
}
