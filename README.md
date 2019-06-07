# Mock Exchange

The mock exchange contains a set of mock companies, that all offer mock stocks. This service creates price fluctuations in these stocks. People can create accounts in order to buy and sell stock, and watch their portfolio grow. 

This service will stream the prices, and provide REST endpoints for user, portfolio, and transaction data.

The web socket streams from `:8080/mockex`

The client that subscribes to this service can be found here: [Mock Exchange Blotter](https://git.csnzoo.com/mduguay/mockex-blotter)

## To run

Build the docker image and run it

`$ docker build -t mock-exchange .`

`$ docker run -p 8000:8080 -d mock-exchange`

## Future

- generate random prices
- stream those prices
- make it containerized
- connect to Postgres db
- create demo user
- host REST endpoints for transactions
- create multiple companies
- allow multiple users/portfolios
