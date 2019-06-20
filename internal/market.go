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
	quotemap := m.getStartQuotes()

	sscan := new(StockScanner)
	stocks := m.Storage.readMultiple(sscan)

	for _, s := range stocks {
		stock, ok := s.(*Stock)
		if !ok {
			log.Println("Error casting stock")
		}
		stock.price = quotemap[stock.symbol].Price
		go stock.startTicking(broadcast)
	}
}

func (s *Stock) startTicking(qPub chan []byte) {
	interval := rand.Intn(1500) + 500
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for range ticker.C {
		s.price = s.tickPrice(s.price)
		q := &Quote{
			Symbol: s.symbol,
			Price:  s.price,
		}
		qbytes, err := json.Marshal(q)
		check(err)
		qPub <- qbytes
	}
}

func (m *Market) getStartQuotes() map[string]*Quote {
	qm := make(map[string]*Quote)
	qscan := new(QuoteScanner)
	quotes := m.Storage.readMultiple(qscan)
	for _, q := range quotes {
		quote, ok := q.(*Quote)
		if !ok {
			log.Println("Error casting quote")
		}
		qm[quote.Symbol] = quote
	}
	return qm
}
