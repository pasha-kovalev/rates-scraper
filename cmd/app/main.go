package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"rates-scraper/internal/config"
	"rates-scraper/internal/controller"
	"rates-scraper/internal/repo"
	"rates-scraper/internal/service"
	"syscall"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Unable to logger", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Unable to load config", err)
	}

	var db *sql.DB
	db, err = repo.NewDb(cfg)
	if err != nil {
		log.Fatal("Unable to connect to db", err)
		return
	}

	ratesRepo := repo.NewRatesRepo(logger, db)

	ratesSvc := service.NewRatesSvc(appCtx, logger, ratesRepo)

	handler := controller.NewHandler(controller.NewRateController(ratesSvc))
	if err = handler.Start(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)); err != nil {
		log.Fatal("Unable to start server", err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	logger.Info("received signal, shutting down", zap.String("signal", sig.String()))
}
