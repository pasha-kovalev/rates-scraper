package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"rates-scraper/internal/service"
	"time"
)

type RateController struct {
	rateService service.RatesService
}

func NewRateController(rateService service.RatesService) *RateController {
	return &RateController{rateService: rateService}
}

func (rc *RateController) GetAllRates(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	rates, err := rc.rateService.GetAllRates(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rates)
}

func (rc *RateController) GetRatesByDate(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()
	date := c.Param(dateParam)
	dateParsed, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return &service.ApiError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	rates, err := rc.rateService.GetOrCollectRatesByDate(ctx, dateParsed)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rates)
}
