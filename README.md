# Mock Exchange

The mock exchange contains a set of mock companies, that all offer mock stocks. This service creates price fluctuations in these stocks. People can create accounts in order to buy and sell stock, and watch their portfolio grow. 

This service will stream the prices, and provide REST endpoints for user, portfolio, and transaction data.

The client that subscribes to this service can be found here: [Mock Exchange Blotter](https://git.csnzoo.com/mduguay/mockex-blotter)

## To run

`go run main.go`

## Future

This will be containerized. This will use a postgres db to store portfolios and company data.
