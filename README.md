application for collecting trade statistics

docker:

```
docker run -d --name _name of docker container_ -p 9000:9000 -e CLICKHOUSE_DB=_your_database_ -e CLICKHOUSE_USER=_your_username_ -e CLICKHOUSE_PASSWORD=_yout_password_ -e
```

DB migration: https://github.com/golang-migrate/migrate

```
command:
migrate create -ext sql -dir migrations/clickhouse init
migrate -path migrations/clickhouse -database "clickhouse://host:port?username=&password=&database=&x-multi-statement=true" up
```

.env:

```
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_USERNAME=
DB_DBNAME=
```
