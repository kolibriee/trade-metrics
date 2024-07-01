CREATE TABLE IF NOT EXISTS order_book
(
    id        UInt32,
    exchange  String,
    pair      String,
    asks      Array(Tuple(Float64, Float64)),
    bids      Array(Tuple(Float64, Float64))
) ENGINE = MergeTree()
ORDER BY (exchange, pair);

CREATE TABLE IF NOT EXISTS order_history 
(
    client_name            String,
    exchange_name          String,
    label                  String,
    pair                   String,
    side                   String,
    type                   String,
    base_qty               Float64,
    price                  Float64,
    algorithm_name_placed  String,
    lowest_sell_prc        Float64,
    highest_buy_prc        Float64,
    commission_quote_qty   Float64,
    time_placed            DateTime
) ENGINE = MergeTree()
ORDER BY (client_name, exchange_name, label, pair);

