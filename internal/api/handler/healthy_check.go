package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterHealthyCheck(router *echo.Echo) {
	// Health check endpoint
	router.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
	})
}
