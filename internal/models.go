package internal

import "time"

type Holding struct {
	Tid    int    `json:"tid"`
	Cid    int    `json:"cid"`
	Symbol string `json:"symbol"`
	Shares int    `json:"shares"`
}

type Trader struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type Trade struct {
	Tid    int     `json:"tid"`
	Cid    int     `json:"cid"`
	Shares int     `json:"shares"`
	Price  float64 `json:"price"`
}

type Cash struct {
	Tid    int     `json:"tid"`
	Amount float64 `json:"amount"`
}

type TradeResult struct {
	Shares int     `json:"shares"`
	Amount float64 `json:"amount"`
}

type HistPoint struct {
	Price float64   `json:"price"`
	Stamp time.Time `json:"stamp"`
}

type Settings struct {
	Cid       int
	Price     float64
	Vol       float64
	Minchange float64
	Maxchange float64
}
