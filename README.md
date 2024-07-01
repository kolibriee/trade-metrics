docker run -d --name _name of docker container_ -p 9000:9000 -e CLICKHOUSE_DB=_your_database_ -e CLICKHOUSE_USER=_your_username_ -e CLICKHOUSE_PASSWORD=_yout_password_ -e CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1 -e CLICKHOUSE_INIT_DEFAULTS=1 clickhouse/clickhouse-server

migrate create -ext sql -dir migrations/clickhouse init
