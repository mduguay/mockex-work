# Mock Exchange

The mock exchange contains a set of mock companies, that all offer mock stocks. This service creates price fluctuations in these stocks. People can create accounts in order to buy and sell stock, and watch their portfolio grow. 

This service will stream the prices, and provide REST endpoints for user, portfolio, and transaction data.

The web socket streams from `:8080/mockex`

## To run

Postgres
```
pg_ctl -D /usr/local/var/postgres start
psql -U postgres
\c mockex
\i scripts/initdb.sql
\dt
\d holding
```

Go
```
go run cmd/mockex/main.go
```

Docker
```
docker build -t mock-exchange .
docker run -p 8080:8080 -d mock-exchange
```

Listen to the feed
```
npm install -g wscat
wscat -c ws://localhost:8080/mockex
```

## Future

- handle incorrect or undefined tid in /holdings handler

- Avoid select when saving quote
- This will be handled when the client sends the id instead of symbol

- Start and stop market via API
 - A kill signal needs to be sent to all stocks

- Make read multiple SQL async - Stress test
- Transaction wrapper for CreateTrade

- Seed historic info
- Historic info API

- Login

- Docker compose Postgres container

- Heroku