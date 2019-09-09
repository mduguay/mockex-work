package internal

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Quote represents an individual tick of a stock
type Quote struct {
	Cid       int       `json:"cid"`
	Timestamp time.Time `json:"timestamp"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Open      float64   `json:"open"`
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
}

// Stock represents a single stock, and associated settings
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
			s.streamTick(qPub)
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

func (s *Stock) backfillTicks(quotechan chan *Quote, timestamp time.Time) {
	stamp := timestamp.Add(interval())
	if stamp.After(time.Now()) {
		close(quotechan)
	}
	s.tickPrice()
	q := s.createQuote(stamp)
	quotechan <- q
}

func (s *Stock) streamTick(qPub chan *Quote) {
	time.Sleep(interval())
	s.tickPrice()
	q := s.createQuote(time.Now())
	qPub <- q
}

func (s *Stock) createQuote(stamp time.Time) *Quote {
	return &Quote{
		Cid:       s.cid,
		Timestamp: stamp,
		Symbol:    s.symbol,
		Price:     s.price,
		Open:      s.price + 3,
		Close:     s.price - 3,
		High:      s.price + 5,
		Low:       s.price - 5,
	}
}

func interval() time.Duration {
	interval := rand.Intn(4000) + 2750
	return time.Duration(interval) * time.Millisecond
}
