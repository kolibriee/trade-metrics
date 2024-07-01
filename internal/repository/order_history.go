package repository

import (
	"context"
	"errors"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kolibriee/trade-metrics/internal/domain"
)

type orderHistoryCH struct {
	db driver.Conn
}

func NewOrderHistoryCH(db driver.Conn) *orderHistoryCH {
	return &orderHistoryCH{
		db: db,
	}
}

func (o *orderHistoryCH) GetOrderHistory(client *domain.Client) ([]*domain.HistoryOrder, error) {
	query := `SELECT client_name, exchange_name, label, pair, side, type,
        base_qty, price, algorithm_name_placed,
        lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed
        FROM order_history
        WHERE client_name = ? AND exchange_name = ? AND label = ? AND pair = ?`

	var orders []*domain.HistoryOrder

	rows, err := o.db.Query(context.Background(), query, client.ClientName, client.ExchangeName, client.Label, client.Pair)
	if err != nil {
		return nil, errors.New("failed to get order history: " + err.Error())
	}
	defer rows.Close()
	var order domain.HistoryOrder
	for rows.Next() {
		err := rows.Scan(
			&order.Client.ClientName,
			&order.Client.ExchangeName,
			&order.Client.Label,
			&order.Client.Pair,
			&order.Side,
			&order.Type,
			&order.BaseQty,
			&order.Price,
			&order.AlgorithmNamePlaced,
			&order.LowestSellPrice,
			&order.HighestBuyPrice,
			&order.CommissionQuoteQty,
			&order.TimePlaced,
		)
		if err != nil {
			return nil, errors.New("failed to scan row: " + err.Error())
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (o *orderHistoryCH) SaveOrder(order *domain.HistoryOrder) error {
	quary := `INSERT INTO order_history (
		client_name, exchange_name, label, pair, side, type,
		base_qty, price, algorithm_name_placed,
		lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if err := o.db.Exec(context.Background(), quary,
		order.Client.ClientName, order.Client.ExchangeName, order.Client.Label, order.Client.Pair,
		order.Side, order.Type, order.BaseQty, order.Price, order.AlgorithmNamePlaced,
		order.LowestSellPrice, order.HighestBuyPrice, order.CommissionQuoteQty,
		order.TimePlaced); err != nil {
		return errors.New("failed to save order: " + err.Error())
	}
	return nil
}
