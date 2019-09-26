package internal

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mduguay/mockex-work/data"
)

// Market is the object that manages all stocks
type Market struct {
	stocks  []*data.Stock
	Storage *Storage
}

// OpeningBell will tell the market to start ticking all stocks
func (m *Market) OpeningBell(broadcast chan []byte) {
	log.Println("Market: Opening Bell")

	m.backfill()

	quotemap := m.getStartQuotes()
	stocks := m.scanStocks()
	stocktick := make(chan *data.Quote)

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

		stocktick := make(chan *data.Quote)
		go stock.backfillTicks(stocktick, starttime)

		for quote := range stocktick {
			// Quote in storage should not have high, low, open, close
			m.Storage.createQuote(quote)
		}
	}
}

func (m *Market) getStartQuotes() map[string]*data.Quote {
	qm := make(map[string]*data.Quote)
	qscan := new(data.QuoteScanner)
	quotes := m.Storage.readMultiple(qscan)
	for _, q := range quotes {
		quote, ok := q.(*data.Quote)
		if !ok {
			log.Println("Error casting quote:", q)
		}
		qm[quote.Symbol] = quote
	}
	return qm
}

func (m *Market) scanStocks() []*data.Stock {
	sscan := new(data.StockScanner)
	rawstocks := m.Storage.readMultiple(sscan)
	var stocks []*data.Stock

	for _, s := range rawstocks {
		stock, ok := s.(*data.Stock)
		if !ok {
			log.Println("Error casting stock:", s)
			continue
		}
		stocks = append(stocks, stock)
	}
	return stocks
}
