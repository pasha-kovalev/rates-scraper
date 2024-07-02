package repo

import (
	"database/sql"
	"fmt"
	"rates-scraper/internal/config"
)

const driveName = "mysql"

func NewDb(config config.Config) (*sql.DB, error) {
	dbCfg := config.Db
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name)

	db, err := sql.Open(driveName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
