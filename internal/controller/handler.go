package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const dateParam = "date"

func NewHandler(controller *RateController) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = ErrorHandler

	e.GET("/rates", controller.GetAllRates)
	e.GET(fmt.Sprintf("/rates/:%s", dateParam), controller.GetRatesByDate)
	return e
}
