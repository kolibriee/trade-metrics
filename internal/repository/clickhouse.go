package repository

import (
	"context"
	"errors"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kolibriee/trade-metrics/internal/config"
)

func NewClickHouseDB(cfg *config.ClickHouse) (driver.Conn, error) {
	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Host + ":" + cfg.Port},
		Auth: clickhouse.Auth{
			Database: cfg.DBName,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return nil, errors.New("failed to connect to ClickHouse: " + err.Error())
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, errors.New("failed to ping ClickHouse: " + err.Error())
	}
	return db, nil
}
