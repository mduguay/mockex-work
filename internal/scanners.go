package internal

import (
	"database/sql"
	"fmt"
	"time"
)

// Scanner will iterate through the rows of the given query and process them
type Scanner interface {
	Query() string
	ScanRow(rows *sql.Rows) interface{}
}

// SingleScanner handles a query that results in a single row
type SingleScanner interface {
	Query() string
	ScanRow(row *sql.Row) interface{}
}

// --- Multi Scanners ---

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

// ScanRow reads a row from the db and creates a Quote
func (qs *QuoteScanner) ScanRow(rows *sql.Rows) interface{} {
	q := new(Quote)
	err := rows.Scan(&q.Cid, &q.Symbol, &q.Price, &q.Timestamp)
	check(err)
	return q
}

// HoldingScanner is responsible for fetching the holdings for a given trader
type HoldingScanner struct {
	uid string
}

// Query is the db query to be executed
func (hs *HoldingScanner) Query() string {
	return fmt.Sprintf(`
		select h.trader_id, c.id, c.symbol, h.shares
		from holding h
		left join company c on h.company_id = c.id
		where h.trader_id = %v`, hs.uid)
}

// ScanRow reads from the db results and creates a Holding
func (hs *HoldingScanner) ScanRow(rows *sql.Rows) interface{} {
	h := new(Holding)
	err := rows.Scan(&h.Tid, &h.Cid, &h.Symbol, &h.Shares)
	check(err)
	return h
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
func (ss *StockScanner) ScanRow(rows *sql.Rows) interface{} {
	s := new(Stock)
	err := rows.Scan(&s.cid, &s.symbol, &s.price, &s.vol, &s.min, &s.max)
	check(err)
	return s
}

// HistoryScanner is responsible for fetching the history of a given stock
type HistoryScanner struct {
	cid int
}

// Query is the db query to be executed
func (hs *HistoryScanner) Query() string {
	t := time.Now()
	location, err := time.LoadLocation("America/New_York")
	if check(err) {
		return ""
	}
	opening := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, location)
	fmt.Println(opening)
	return fmt.Sprintf("select price, stamp from quote where company_id = %v and stamp > now()::date order by stamp desc", hs.cid)
}

// ScanRow reads the results from storage and creates a HistPoint
func (hs *HistoryScanner) ScanRow(rows *sql.Rows) interface{} {
	hp := new(HistPoint)
	err := rows.Scan(&hp.Price, &hp.Stamp)
	check(err)
	return hp
}

// --- Single Scanners ---

// TraderScanner is responsible for fetching a single trader given an id
type TraderScanner struct {
	id int
}

// Query is the db query to be executed
func (ts *TraderScanner) Query() string {
	return fmt.Sprintf("select email from trader where id = %v", ts.id)
}

// ScanRow handles the result given from the db query above
func (ts *TraderScanner) ScanRow(stmt *sql.Stmt) interface{} {
	t := new(Trader)
	err := stmt.QueryRow(ts.id).Scan(t.Email)
	check(err)
	return t
}

// CashScanner is responsible for reading the cash value from the db for the given trader
type CashScanner struct {
	id int
}

// Query is the db query to be executed
func (cs *CashScanner) Query() string {
	return fmt.Sprintf("select amount from cash where trader_id = %v", cs.id)
}

// ScanRow handles the result given from the db query above
func (cs *CashScanner) ScanRow(stmt *sql.Stmt) interface{} {
	c := new(Cash)
	err := stmt.QueryRow(cs.id).Scan(c.Amount)
	check(err)
	return c
}
