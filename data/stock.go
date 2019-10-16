package data

import (
	"database/sql"
	"log"
	"math"
	"math/rand"
	"time"
)

// Stock represents a single stock, and associated settings
type Stock struct {
	cid         int
	Symbol      string
	Price       float64
	min         float64
	max         float64
	vol         float64
	Stopchan    chan struct{}
	Stoppedchan chan struct{}
}

// StockScanner is responsible for fetching stocks form the db for a given company
type StockScanner struct{}

// Query is the db query to be executed
func (ss *StockScanner) Query() string {
	return `
		select c.id, c.symbol, s.vol, s.minchange, s.maxchange
		from stock s
		left join company c
		on s.company_id = c.id
		`
}

// ScanRow reads the results from storage and creates a Stock
func (ss *StockScanner) ScanRow(rows *sql.Rows) (interface{}, error) {
	s := new(Stock)
	err := rows.Scan(&s.cid, &s.Symbol, &s.vol, &s.min, &s.max)
	return s, err
}

// GenerateTicks will create and stream stock ticks to the qPub channel
func (s *Stock) GenerateTicks(qPub chan *Quote) {
	defer close(s.Stoppedchan)
	for {
		select {
		default:
			s.streamTick(qPub)
		case <-s.Stopchan:
			log.Println("Stopping tick generation")
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
	changeAmt := s.Price * changePct
	newPrice := s.Price + changeAmt
	rPrice := math.Round(newPrice*100) / 100
	s.Price = rPrice
}

// BackfillTicks will create quotes from the given stamp until now
func (s *Stock) BackfillTicks(quotechan chan *Quote, stamp time.Time) {
	now := time.Now().UTC()
	for stamp.Before(now) {
		i := interval()
		stamp = stamp.Add(i)
		s.tickPrice()
		q := s.createQuote(stamp)
		quotechan <- q
	}
	close(quotechan)
}

func (s *Stock) streamTick(qPub chan *Quote) {
	time.Sleep(interval())
	s.tickPrice()
	q := s.createQuote(time.Now().UTC())
	qPub <- q
}

func (s *Stock) createQuote(stamp time.Time) *Quote {
	return &Quote{
		Cid:       s.cid,
		Timestamp: stamp,
		Symbol:    s.Symbol,
		Price:     s.Price,
	}
}

func interval() time.Duration {
	interval := rand.Intn(4000) + 2750
	return time.Duration(interval) * time.Millisecond
}
