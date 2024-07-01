package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kolibriee/trade-metrics/internal/domain"
)

type orderBookCH struct {
	db driver.Conn
}

func NewOrderBookCH(db driver.Conn) *orderBookCH {
	return &orderBookCH{
		db: db,
	}
}

func (o *orderBookCH) GetOrderBook(exchangeName, pair string) (*domain.AsksBids, error) {
	query := `
        SELECT id, asks, bids 
        FROM order_book FINAL
        WHERE exchange = ? AND pair = ?
    `

	var (
		id   uint32
		asks [][]float64
		bids [][]float64
	)

	err := o.db.QueryRow(context.Background(), query, exchangeName, pair).Scan(&id, &asks, &bids)
	if err != nil {
		return nil, errors.New("failed to get order book: " + err.Error())
	}

	var asksBids domain.AsksBids
	asksBids.Id = id
	asksBids.Asks = make([]domain.DepthOrder, len(asks))
	for i, ask := range asks {
		asksBids.Asks[i] = domain.DepthOrder{
			Price:   ask[0],
			BaseQty: ask[1],
		}
	}
	asksBids.Bids = make([]domain.DepthOrder, len(bids))
	for i, bid := range bids {
		asksBids.Bids[i] = domain.DepthOrder{
			Price:   bid[0],
			BaseQty: bid[1],
		}
	}

	return &asksBids, nil
}

func (o *orderBookCH) SaveOrderBook(exchangeName, pair string, asksBids *domain.AsksBids) error {
	id := asksBids.Id
	asks := make([]string, len(asksBids.Asks))
	for i, ask := range asksBids.Asks {
		asks[i] = fmt.Sprintf("(%f,%f)", ask.Price, ask.BaseQty)
	}

	bids := make([]string, len(asksBids.Bids))
	for i, bid := range asksBids.Bids {
		bids[i] = fmt.Sprintf("(%f,%f)", bid.Price, bid.BaseQty)
	}

	query := `
        INSERT INTO order_book (id, exchange, pair, asks, bids)
        VALUES (?, ?, ?, ?, ?)
    `
	if err := o.db.Exec(context.Background(), query, id, exchangeName, pair, asks, bids); err != nil {
		return errors.New("failed to save order book: " + err.Error())
	}
	return nil
}
