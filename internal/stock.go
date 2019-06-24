package internal

import (
	"math"
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

func (s *Stock) generateTicks(qPub chan *Quote) {
	for {
		interval := rand.Intn(2000) + 750
		//interval := 0
		time.Sleep(time.Duration(interval) * time.Millisecond)
		s.tickPrice()
		q := &Quote{
			Symbol: s.symbol,
			Price:  s.price,
		}
		qPub <- q
	}
}

func (s *Stock) tickPrice() {
	rnd := s.min + rand.Float64()*(s.max-s.min)
	changePct := 2 * s.vol * rnd
	if changePct > s.vol {
		changePct -= (2 * s.vol)
	}
	changeAmt := s.price * changePct
	newPrice := s.price + changeAmt
	rPrice := math.Round(newPrice*100) / 100
	s.price = rPrice
}
