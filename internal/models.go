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
	Shares    int     `json:"shares"`
	Direction bool    `json:"direction"`
	Price     float64 `json:"price"`
}
