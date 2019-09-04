package internal

import "time"

// Holding represents a traders assets in a particular stock
type Holding struct {
	Tid    int    `json:"tid"`
	Cid    int    `json:"cid"`
	Symbol string `json:"symbol"`
	Shares int    `json:"shares"`
}

// Trader is a user that holds stock and makes trades
type Trader struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// Trade is an action taken by a trader to buy or sell stocks
type Trade struct {
	Tid    int     `json:"tid"`
	Cid    int     `json:"cid"`
	Shares int     `json:"shares"`
	Price  float64 `json:"price"`
}

// Cash is the amount of money on hand that can be used to buy stocks
type Cash struct {
	Tid    int     `json:"tid"`
	Amount float64 `json:"amount"`
}

// TradeResult captures the results of a Trade action and shows the resulting holdings
type TradeResult struct {
	Shares int     `json:"shares"`
	Amount float64 `json:"amount"`
}

// HistPoint is one point in the history of a stock's price
type HistPoint struct {
	Price float64   `json:"price"`
	Stamp time.Time `json:"stamp"`
}

// Settings holds the values that define how a stock is priced and how that price changes
type Settings struct {
	Cid       int
	Price     float64
	Vol       float64
	Minchange float64
	Maxchange float64
}
