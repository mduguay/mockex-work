package internal

import (
	"encoding/json"
	"log"
	"time"
)

// Market is the object that manages all stocks
type Market struct {
	stocks  []*Stock
	Storage *Storage
}

// OpeningBell will tell the market to start ticking all stocks
func (m *Market) OpeningBell(broadcast chan []byte) {
	log.Println("Market: Opening Bell")

	m.backfill()

	quotemap := m.getStartQuotes()
	stocks := m.scanStocks()
	stocktick := make(chan *Quote)

	for _, stock := range stocks {
		stopchan := make(chan struct{})
		stoppedchan := make(chan struct{})
		stock.stopchan = stopchan
		stock.stoppedchan = stoppedchan
		stock.price = quotemap[stock.symbol].Price
		m.stocks = append(m.stocks, stock)
		go stock.generateTicks(stocktick)
	}

	for quote := range stocktick {
		m.Storage.createQuote(quote)
		qbytes, err := json.Marshal(quote)
		check(err)
		broadcast <- qbytes
	}
}

// ClosingBell will stop the stocks from ticking
func (m *Market) ClosingBell() {
	log.Println("Market: Closing Bell")
	for _, s := range m.stocks {
		close(s.stopchan)
		<-s.stoppedchan
	}
}

func (m *Market) backfill() {
	log.Println("Market: Backfilling")
	quotemap := m.getStartQuotes()
	stocks := m.scanStocks()

	for _, stock := range stocks {
		stock.price = quotemap[stock.symbol].Price

		now := time.Now()
		todayopen := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC)
		starttime := quotemap[stock.symbol].Timestamp
		if starttime.Before(todayopen) {
			starttime = todayopen
		}

		stocktick := make(chan *Quote)
		go stock.backfillTicks(stocktick, starttime)

		for quote := range stocktick {
			// Quote in storage should not have high, low, open, close
			m.Storage.createQuote(quote)
		}
	}
}

func (m *Market) getStartQuotes() map[string]*Quote {
	qm := make(map[string]*Quote)
	qscan := new(QuoteScanner)
	quotes := m.Storage.readMultiple(qscan)
	for _, q := range quotes {
		quote, ok := q.(*Quote)
		if !ok {
			log.Println("Error casting quote:", q)
		}
		qm[quote.Symbol] = quote
	}
	return qm
}

func (m *Market) scanStocks() []*Stock {
	sscan := new(StockScanner)
	rawstocks := m.Storage.readMultiple(sscan)
	var stocks []*Stock

	for _, s := range rawstocks {
		stock, ok := s.(*Stock)
		if !ok {
			log.Println("Error casting stock:", s)
			continue
		}
		stocks = append(stocks, stock)
	}
	return stocks
}

func (m *Market) history() []*Quote {
	// scan all quotes in time interval
	// for now, just today
	qscan := new(QuoteScanner)
	quotes := m.Storage.readMultiple(qscan)
	// figure out increment amount (5 mins)
	// sort by time
	// loop through all quotes in increment
	// record start, end, high, low (sort for high/low)
	// stamp quote with start time or end time?
	return &ChartPoint{
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
