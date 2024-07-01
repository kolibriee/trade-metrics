package domain

type AsksBids struct {
	Id   uint32       `db:"id" json:"id"`
	Asks []DepthOrder `json:"asks" binding:"required"`
	Bids []DepthOrder `json:"bids" binding:"required"`
}

type OrderBook struct {
	ID       int64        `db:"id"`
	Exchange string       `db:"exchange"`
	Pair     string       `db:"pair"`
	Asks     []DepthOrder `db:"asks"`
	Bids     []DepthOrder `db:"bids"`
}

type DepthOrder struct {
	Price   float64 `db:"price" json:"price" binding:"required"`
	BaseQty float64 `db:"base_qty" json:"base_qty" binding:"required"`
}
