docker run -d --name trade-metrics-clickhouse -p 9000:9000 -e CLICKHOUSE_DB=trade_collect -e CLICKHOUSE_USER=kolibriee -e CLICKHOUSE_PASSWORD=1w2qaxsz3EDC -e CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1 -e CLICKHOUSE_INIT_DEFAULTS=1 clickhouse/clickhouse-server

migrate create -ext sql -dir migrations/clickhouse init
