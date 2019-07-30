package internal

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Quote struct {
	Cid    int     `json:"cid"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Stock struct {
	cid         int
	symbol      string
	price       float64
	min         float64
	max         float64
	vol         float64
	stopchan    chan struct{}
	stoppedchan chan struct{}
}

func (s *Stock) generateTicks(qPub chan *Quote) {
	defer close(s.stoppedchan)
	for {
		select {
		default:
			interval := rand.Intn(2000) + 750
			//interval := 0
			time.Sleep(time.Duration(interval) * time.Millisecond)
			s.tickPrice()
			q := &Quote{
				Cid:    s.cid,
				Symbol: s.symbol,
				Price:  s.price,
			}
			qPub <- q
		case <-s.stopchan:
			fmt.Println("Stock.stopchan")
			return
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
	newPrice := s.price + changeAmt
	rPrice := math.Round(newPrice*100) / 100
	s.price = rPrice
}
