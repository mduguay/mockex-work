package internal

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

type Company struct {
	Id     int
	Symbol string
}

type Quote struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Stock struct {
	symbol string
	price  float64
	min    float64
	max    float64
	vol    float64
}

type Holding struct {
	Uid    int
	Symbol string
	Shares int
}

type Market struct {
	stocks  []*Stock
	Storage *Storage
}

const (
	start = 30.0
	min   = 0.0
	max   = 1.0
	vol   = 0.02
)

func genNextPrice(oldPrice float64) float64 {
	rnd := min + rand.Float64()*(max-min)
	changePct := 2 * vol * rnd
	if changePct > vol {
		changePct -= (2 * vol)
	}
	changeAmt := oldPrice * changePct
	newPrice := oldPrice + changeAmt
	return newPrice
}

func getDayPrices() []float64 {
	prices := make([]float64, 96)
	prices[0] = start
	for i := 1; i < 96; i++ {
		prices[i] = genNextPrice(prices[i-1])
	}
	return prices
}

func (m *Market) Prime() {
	cs := new(CompanyScanner)
	companies := m.Storage.readMultiple(cs)

	m.stocks = make([]*Stock, len(companies))
	for i, c := range companies {
		m.stocks[i] = &Stock{
			symbol: c.(*Company).Symbol,
			price:  30.00,
			min:    0.00,
			max:    1.00,
			vol:    0.02,
		}
	}
}

func (s *Stock) tickPrice() {
	s.price = genNextPrice(s.price)
}

func (m *Market) OpeningBell(broadcast chan []byte) {
	qs := make([]*Quote, len(m.stocks))
	for {
		time.Sleep(time.Second * 2)
		for i, s := range m.stocks {
			s.tickPrice()
			q := &Quote{
				Symbol: s.symbol,
				Price:  s.price,
			}
			qs[i] = q
		}
		qsb, err := json.Marshal(qs)
		if err != nil {
			log.Println(err)
		}
		broadcast <- qsb
	}
}

// func generateHistoricData() {
// 	p := getDayPrices()
// 	//insert into db
// }
