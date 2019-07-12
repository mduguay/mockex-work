package internal

type Holding struct {
	Tid    int    `json:"tid"`
	Symbol string `json:"symbol"`
	Shares int    `json:"shares"`
}

type Trader struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type Trade struct {
	Tid       int     `json:"tid"`
	Cid       int     `json:"cid"`
	Shares    int     `json:"shares"`
	Direction bool    `json:"direction"`
	Price     float64 `json:"price"`
}
