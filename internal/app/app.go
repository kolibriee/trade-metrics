package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kolibriee/trade-metrics/internal/config"
	"github.com/kolibriee/trade-metrics/internal/repository"
	"github.com/kolibriee/trade-metrics/internal/server"
	"github.com/sirupsen/logrus"

	"github.com/kolibriee/trade-metrics/internal/controller"
)

func Run(configsDir string, configName string) {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	config, err := config.New(configsDir, configName)
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := repository.NewClickHouseDB(&config.ClickHouse)
	if err != nil {
		logrus.Fatalf("failed to connect to ClickHouse: %v", err)
	}

	repo := repository.NewRepository(db)
	controller := controller.NewController(repo)
	var srv server.Server
	go func() {
		if err := srv.Run(&config.Server, controller.Handler); err != nil {
			logrus.Fatalf("failed to start server: %v", err)
		}
	}()
	logrus.Info("App started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
