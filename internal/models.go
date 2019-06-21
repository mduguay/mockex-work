package internal

type Holding struct {
	Uid    int
	Symbol string
	Shares int
}

type Trader struct {
	Id    int
	Email string
}

type Trade struct {
	Tid       int     `json:"tid"`
	Sid       int     `json:"sid"`
	Amount    int     `json:"amount"`
	Direction string  `json:"direction"`
	Price     float64 `json:"price"`
}
