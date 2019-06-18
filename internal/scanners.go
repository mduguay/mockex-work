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

func (cs *CompanyScanner) ScanRow(rows *sql.Rows) interface{} {
	c := new(Company)
	err := rows.Scan(&c.Id, &c.Symbol)
	check(err)
	return c
}

// Should return | symbol | lastprice |
func (cs *CompanyScanner) Query() string {
	return `
		select c.symbol, p.price
		from company c
		left join ( select p.company_id, price p
		on c.id = p.company_id
		order by p.stamp desc
		group by c.symbol;
		`
}

type HoldingScanner struct {
	uid string
}

func (hs *HoldingScanner) ScanRow(rows *sql.Rows) interface{} {
	h := new(Holding)
	err := rows.Scan(&h.Uid, &h.Symbol, &h.Shares)
	check(err)
	return h
}

func (hs *HoldingScanner) Query() string {
	return fmt.Sprintf(`
		select h.trader_id, c.symbol, h.shares
		from holding h
		left join company c on h.company_id = c.id
		where h.trader_id = %v`, hs.uid)
}
