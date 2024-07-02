package repository

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kolibriee/trade-metrics/internal/domain"
)

type Orderbook interface {
	GetOrderBook(exchangeName, pair string) (*domain.AsksBids, error)
	SaveOrderBook(exchangeName, pair string, asksBids *domain.AsksBids) error
}

type Orderhistory interface {
	GetOrderHistory(client *domain.Client) ([]*domain.HistoryOrder, error)
	SaveOrder(order *domain.HistoryOrder) error
}

type Repository struct {
	Orderbook
	Orderhistory
}

func NewRepository(db driver.Conn) *Repository {
	return &Repository{
		Orderbook:    NewOrderBookCH(db),
		Orderhistory: NewOrderHistoryCH(db),
	}
}
