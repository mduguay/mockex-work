package internal

import (
	"encoding/json"
	"fmt"
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
	m.backfill()
	log.Println("Market: Opening Bell")
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
	fmt.Println("Market.ClosingBell")
	for _, s := range m.stocks {
		fmt.Println("Stopping ", s.symbol)
		close(s.stopchan)
		<-s.stoppedchan
		fmt.Println("Stopped")
	}
}

func (m *Market) backfill() {
	log.Println("Market: Backfilling")
	quotemap := m.getStartQuotes()
	stocks := m.scanStocks()

	for _, stock := range stocks {
		stock.price = quotemap[stock.symbol].Price

		now := time.Now()
		todayopen := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.UTC)
		starttime := quotemap[stock.symbol].Timestamp
		if starttime.Before(todayopen) {
			starttime = todayopen
		}

		stocktick := make(chan *Quote)
		go stock.backfillTicks(stocktick, starttime)

		for quote := range stocktick {
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
