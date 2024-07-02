package repo

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const (
	ratesDateLayout = "2006-01-02"
	ratesDateFormat = "%04d-%02d-%02d"
	dbDateLayout    = `"2006-01-02 15:04:05"`
)
const (
	getAllQuery    = "SELECT id, cur_id, rate_date, cur_abbreviation, cur_scale, cur_name, cur_official_rate FROM rates "
	getByDateQuery = getAllQuery + "WHERE rate_date >= ? AND rate_date < ?"
	insertQuery    = "INSERT INTO rates (cur_id, rate_date, cur_abbreviation, cur_scale, cur_name, cur_official_rate) VALUES (?, ?, ?, ?, ?, ?)"
)

type RatesRepo interface {
	InsertAll(ctx context.Context, rates []RateEntity) error
	GetAll(ctx context.Context) ([]RateEntity, error)
	GetByDate(ctx context.Context, day int, month int, year int) ([]RateEntity, error)
}

type ratesRepo struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRatesRepo(logger *zap.Logger, db *sql.DB) RatesRepo {
	return &ratesRepo{
		logger: logger,
		db:     db}
}

func (r *ratesRepo) GetAll(ctx context.Context) ([]RateEntity, error) {
	rows, err := r.db.QueryContext(ctx, getAllQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRateEntities(rows)

}

func (r *ratesRepo) GetByDate(ctx context.Context, day int, month int, year int) ([]RateEntity, error) {
	date, err := time.Parse(ratesDateLayout, fmt.Sprintf(ratesDateFormat, year, month, day))
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx,
		getByDateQuery,
		date,
		date.Add(24*time.Hour))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRateEntities(rows)
}

func (r *ratesRepo) InsertAll(ctx context.Context, rates []RateEntity) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, rate := range rates {
		_, err = stmt.ExecContext(ctx,
			rate.CurID,
			time.Time(rate.RateDate).Format(ratesDateLayout),
			rate.CurAbbreviation,
			rate.CurScale,
			rate.CurName,
			rate.CurOfficialRate)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func scanRateEntities(rows *sql.Rows) ([]RateEntity, error) {
	var rates []RateEntity
	for rows.Next() {
		var rate RateEntity
		err := rows.Scan(
			&rate.ID,
			&rate.CurID,
			&rate.RateDate,
			&rate.CurAbbreviation,
			&rate.CurScale,
			&rate.CurName,
			&rate.CurOfficialRate,
		)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}
