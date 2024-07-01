package domain

import "time"

type HistoryOrder struct {
	Client              Client    `db:"client" json:"client" binding:"required"`
	Side                string    `db:"side" json:"side" binding:"required"`
	Type                string    `db:"type" json:"type" binding:"required"`
	BaseQty             float64   `db:"base_qty" json:"base_qty" binding:"required"`
	Price               float64   `db:"price" json:"price" binding:"required"`
	AlgorithmNamePlaced string    `db:"algorithm_name_placed" json:"algorithm_name_placed" binding:"required"`
	LowestSellPrice     float64   `db:"lowest_sell_prc" json:"lowest_sell_prc" binding:"required"`
	HighestBuyPrice     float64   `db:"highest_buy_prc" json:"highest_buy_prc" binding:"required"`
	CommissionQuoteQty  float64   `db:"commission_quote_qty" json:"commission_quote_qty" binding:"required"`
	TimePlaced          time.Time `db:"time_placed" json:"time_placed"`
}

type Client struct {
	ClientName   string `db:"client_name" json:"client_name" binding:"required"`
	ExchangeName string `db:"exchange_name" json:"exchange_name" binding:"required"`
	Label        string `db:"label" json:"label" binding:"required"`
	Pair         string `db:"pair" json:"pair" binding:"required"`
}
