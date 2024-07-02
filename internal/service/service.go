package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"rates-scraper/internal/repo"
	"time"
)

const (
	ratesApiBaseUrl = "https://api.nbrb.by/exrates/rates?periodicity=0"
	schedulerPeriod = 24 * time.Hour
)

type RatesService interface {
	StartScheduler() error
	GetAllRates(ctx context.Context) (*RateDto, error)
	GetOrCollectRatesByDate(ctx context.Context, date time.Time) (*RateDto, error)
	collectRates(ctx context.Context, date *time.Time) (*RateDto, error)
}

type ratesService struct {
	appCtx      context.Context
	logger      *zap.Logger
	ratesRepo   repo.RatesRepo
	ratesApiUrl string
}

func NewRatesSvc(appCtx context.Context, logger *zap.Logger, ratesRepo repo.RatesRepo) RatesService {
	return &ratesService{
		logger:    logger,
		ratesRepo: ratesRepo,
		appCtx:    appCtx,
	}
}

func (s *ratesService) StartScheduler() error {
	ctx := context.Background()
	now := time.Now()

	_, err := s.GetOrCollectRatesByDate(ctx, now)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(schedulerPeriod)
	go func() {
		for {
			select {
			case <-s.appCtx.Done():
				return
			case <-ticker.C:
				if _, err = s.GetOrCollectRatesByDate(s.appCtx, now); err != nil {
					s.logger.Error("Unable to collect rates", zap.Error(err))
				}
			}
		}
	}()
	return nil
}

func (s *ratesService) GetAllRates(ctx context.Context) (*RateDto, error) {
	rates, err := s.ratesRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &RateDto{rates}, nil
}

func (s *ratesService) GetOrCollectRatesByDate(ctx context.Context, date time.Time) (*RateDto, error) {
	rates, err := s.ratesRepo.GetByDate(ctx, date.Day(), int(date.Month()), date.Year())
	if errors.Is(err, sql.ErrNoRows) || len(rates) == 0 {
		return s.collectRates(ctx, &date)
	}
	if err != nil {
		return nil, err
	}

	return &RateDto{rates}, nil
}

func (s *ratesService) collectRates(ctx context.Context, date *time.Time) (*RateDto, error) {
	url := ratesApiBaseUrl
	if date != nil {
		url = makeGetByDateUrl(*date)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrBadGateway
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("received response", zap.String("body", string(bodyBytes)))
	var data []repo.RateEntity
	if err = json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}

	return &RateDto{ExchangeRates: data}, s.ratesRepo.InsertAll(ctx, data)
}

func makeGetByDateUrl(date time.Time) string {
	return fmt.Sprintf("%s&ondate=%s", ratesApiBaseUrl, date.Format(time.DateOnly))
}
