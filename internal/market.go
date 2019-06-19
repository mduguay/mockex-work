package internal

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

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

// func getDayPrices() []float64 {
// 	prices := make([]float64, 96)
// 	prices[0] = start
// 	for i := 1; i < 96; i++ {
// 		prices[i] = genNextPrice(prices[i-1])
// 	}
// 	return prices
// }

func (m *Market) Prime() {
	ss := new(StockScanner)
	stocks := m.Storage.readMultiple(ss)

	m.stocks = make([]*Stock, len(stocks))
	for i, s := range stocks {
		m.stocks[i] = &Stock{
			symbol: s.(*Stock).symbol,
			price:  30.00,
			min:    0.00,
			max:    1.00,
			vol:    0.02,
		}
	}
}

func (s *Stock) tickPrice() {
	rnd := s.min + rand.Float64()*(s.max-s.min)
	changePct := 2 * s.vol * rnd
	if changePct > s.vol {
		changePct -= (2 * s.vol)
	}
	changeAmt := s.price * changePct
	s.price = s.price + changeAmt
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
