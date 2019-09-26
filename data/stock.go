package data

import "database/sql"

// Stock represents a single stock, and associated settings
type Stock struct {
	cid         int
	symbol      string
	price       float64
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
		select c.id, c.symbol, s.price, s.vol, s.minchange, s.maxchange
		from stock s
		left join company c
		on s.company_id = c.id
		`
}

// ScanRow reads the results from storage and creates a Stock
func (ss *StockScanner) ScanRow(rows *sql.Rows) (interface{}, error) {
	s := new(Stock)
	err := rows.Scan(&s.cid, &s.symbol, &s.price, &s.vol, &s.min, &s.max)
	return s, err
}
