# Mock Exchange

The mock exchange contains a set of mock companies, that all offer mock stocks. This service creates price fluctuations in these stocks. People can create accounts in order to buy and sell stock, and watch their portfolio grow. 

This service will stream the prices, and provide REST endpoints for user, portfolio, and transaction data.

The web socket streams from `:8080/mockex`

This is a prime example of goroutines in generating stock prices. Postgres is used in the back end for a datastore.

## To run

### Postgres
```
docker pull postgres
mkdir -p $HOME/docker/volumes/postgres
docker run --rm   --name pg-docker -e POSTGRES_PASSWORD=postgres64 -d -p 5432:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data  postgres
```

To setup the db for the first time (from within psql)
```
psql -h localhost -U postgres -d postgres

create database mockex
\c mockex
\i scripts/initdb.sql
\i scripts/seed.sql
\d
SELECT * FROM company;
```

---

### Environment variables

Make a copy of the dev environment variables
```
cp .env_dev .env
```

Adjust as needed if your settings differ from the defaults

---

### Run the server

Go - run the main application
```
go run cmd/mockex/main.go
```

Or, if you're using vscode (highly recommended), just hit `F5`

---

### Test the API

Listen to the feed
```
npm install -g wscat
wscat -c ws://localhost:8080/mockex
```

Get stock info
```
curl http://localhost:8080/quotes
```

---

### Using Postman

Open up Postman, click "Import" in the top left, then point it to the `Mockex.postman_collection.json` file

---

### Docker
```
docker build -t mock-exchange .
docker run -p 8080:8080 -d mock-exchange
```

---

## Future

- Adjust stock variability via API

- handle incorrect or undefined tid in /holdings handler

- Avoid select when saving quote
- This will be handled when the client sends the id instead of symbol

- Docker compose Postgres container

- Start and stop market via API
- A kill signal needs to be sent to all stocks

- Make read multiple SQL async - Stress test
- Transaction wrapper for CreateTrade

- Seed historic info
- Historic info API

- Login

- Heroku