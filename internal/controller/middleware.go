package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"rates-scraper/internal/service"
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := err.Error()
	var he *service.ApiError
	if errors.As(err, &he) {
		code = he.Code
	}
	c.JSON(code, echo.Map{
		"error": message,
	})
}
