package internal

import (
	"database/sql"
	"fmt"
)

type Scanner interface {
	Query() string
	ScanRow(rows *sql.Rows) interface{}
}

type SingleScanner interface {
	Query() string
	ScanRow(row *sql.Row) interface{}
}

type QuoteScanner struct{}

func (qs *QuoteScanner) Query() string {
	return `
		select c.symbol, q.price lastprice
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

func (qs *QuoteScanner) ScanRow(rows *sql.Rows) interface{} {
	q := new(Quote)
	err := rows.Scan(&q.Symbol, &q.Price)
	check(err)
	return q
}

type HoldingScanner struct {
	uid string
}

func (hs *HoldingScanner) Query() string {
	return fmt.Sprintf(`
		select h.trader_id, c.symbol, h.shares
		from holding h
		left join company c on h.company_id = c.id
		where h.trader_id = %v`, hs.uid)
}

func (hs *HoldingScanner) ScanRow(rows *sql.Rows) interface{} {
	h := new(Holding)
	err := rows.Scan(&h.Uid, &h.Symbol, &h.Shares)
	check(err)
	return h
}

type StockScanner struct{}

func (ss *StockScanner) Query() string {
	return `
		select c.symbol, s.price, s.vol, s.minchange, s.maxchange
		from stock s
		left join company c
		on s.company_id = c.id
		`
}

func (ss *StockScanner) ScanRow(rows *sql.Rows) interface{} {
	s := new(Stock)
	err := rows.Scan(&s.symbol, &s.price, &s.vol, &s.min, &s.max)
	check(err)
	return s
}

type TraderScanner struct {
	id int
}

func (ts *TraderScanner) Query() string {
	return fmt.Sprintf("select email from trader where id = %v", ts.id)
}

func (ss *TraderScanner) ScanRow(stmt *sql.Stmt) interface{} {
	t := new(Trader)
	err := stmt.QueryRow(ss.id).Scan(t.Email)
	check(err)
	return t
}
