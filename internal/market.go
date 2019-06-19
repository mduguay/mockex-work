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
		m.stocks[i] = s.(*Stock)
	}
}

func (s *Stock) tickPrice(p float64) float64 {
	rnd := s.min + rand.Float64()*(s.max-s.min)
	changePct := 2 * s.vol * rnd
	if changePct > s.vol {
		changePct -= (2 * s.vol)
	}
	changeAmt := p * changePct
	return p + changeAmt
}

func (m *Market) OpeningBell(broadcast chan []byte) {
	qs := make([]*Quote, len(m.stocks))
	// Fetch quotes
	// qs := new(QuoteScanner)
	// quotes := m.Storage.readMultiple(qs)
	// There's a 1:1 mapping between quotes (latest price) and stocks (metadata)
	// How should these be linked
	for {
		time.Sleep(time.Second * 2)
		for i, s := range m.stocks {
			s.price = s.tickPrice(s.price)
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
