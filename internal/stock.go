package internal

import (
	"encoding/json"
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

func (s *Stock) startTicking(qPub chan []byte) {
	interval := rand.Intn(2000) + 750
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for range ticker.C {
		s.tickPrice()
		q := &Quote{
			Symbol: s.symbol,
			Price:  s.price,
		}
		qbytes, err := json.Marshal(q)
		check(err)
		qPub <- qbytes
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
