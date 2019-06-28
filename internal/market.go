package internal

import (
	"encoding/json"
	"fmt"
	"log"
)

type Market struct {
	stocks  []*Stock
	Storage *Storage
}

func (m *Market) OpeningBell(broadcast chan []byte) {
	quotemap := m.getStartQuotes()

	sscan := new(StockScanner)
	stocks := m.Storage.readMultiple(sscan)

	stocktick := make(chan *Quote)

	for _, s := range stocks {
		stock, ok := s.(*Stock)
		if !ok {
			log.Println("Error casting stock:", s)
		}
		fmt.Println(quotemap[stock.symbol])
		stock.price = quotemap[stock.symbol].Price
		go stock.generateTicks(stocktick)
	}

	for quote := range stocktick {
		m.Storage.createQuote(quote)
		qbytes, err := json.Marshal(quote)
		check(err)
		broadcast <- qbytes
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
