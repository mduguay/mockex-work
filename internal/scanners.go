package internal

import (
	"database/sql"
	"fmt"
)

type Scanner interface {
	ScanRow(rows *sql.Rows) interface{}
	Query() string
}

type CompanyScanner struct{}

func (cs *CompanyScanner) Query() string {
	return `
		select c.symbol, p.price lastprice
		from price p
		right join (
			select company_id cid, max(stamp) ms
			from price
			group by company_id
			) laststamp
		on p.stamp = laststamp.ms
		and p.company_id = laststamp.cid
		left join company c
		on c.id = p.company_id;
		`
}

func (cs *CompanyScanner) ScanRow(rows *sql.Rows) interface{} {
	c := new(Company)
	err := rows.Scan(&c.Symbol, &c.LastPrice)
	check(err)
	return c
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
