package service

import "rates-scraper/internal/repo"

type RateDto struct {
	ExchangeRates []repo.RateEntity
}
