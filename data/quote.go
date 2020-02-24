package data

import (
	"database/sql"
	"time"
)

// Quote represents an individual tick of a stock
type Quote struct {
	Cid       int       `json:"cid"`
	Timestamp time.Time `json:"timestamp"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
}

// QuoteScanner reads quotes from the database
type QuoteScanner struct{}

// Query is the db query to be executed
func (qs *QuoteScanner) Query() string {
	return `
		select c.id, c.symbol, q.price lastprice, q.stamp
		from quote q
		right join (
			select company_id cid, max(stamp) ms
			from quote
			group by company_id
			) laststamp
		on q.stamp = laststamp.ms
		and q.company_id = laststamp.cid
		left join company c
		on c.id = q.company_id;
		`
}

// Params returns the parameters for the query
func (qs *QuoteScanner) Params() []interface{} {
	return []interface{}{}
}

// ScanRow reads a row from the db and creates a Quote
func (qs *QuoteScanner) ScanRow(rows *sql.Rows) (interface{}, error) {
	q := new(Quote)
	err := rows.Scan(&q.Cid, &q.Symbol, &q.Price, &q.Timestamp)
	return q, err
}
