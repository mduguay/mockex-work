package internal

// Trade is an action taken by a trader to buy or sell stocks
type Trade struct {
	Tid    int     `json:"tid"`
	Cid    int     `json:"cid"`
	Shares int     `json:"shares"`
	Price  float64 `json:"price"`
}

// TradeResult captures the results of a Trade action and shows the resulting holdings
type TradeResult struct {
	Shares int     `json:"shares"`
	Amount float64 `json:"amount"`
}

// Settings holds the values that define how a stock is priced and how that price changes
type Settings struct {
	Cid       int
	Price     float64
	Vol       float64
	Minchange float64
	Maxchange float64
}
