package internal

import (
	"database/sql"
	"fmt"
)

type Scanner interface {
	ScanRow(rows *sql.Rows) interface{}
	GetQuery() string
}

type CompanyScanner struct{}

func (cs *CompanyScanner) ScanRow(rows *sql.Rows) interface{} {
	c := new(Company)
	err := rows.Scan(&c.Id, &c.Symbol)
	check(err)
	return c
}

func (cs *CompanyScanner) GetQuery() string {
	return "select id, symbol from company"
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

func (hs *HoldingScanner) GetQuery() string {
	return fmt.Sprintf(`
		select h.trader_id, c.symbol, h.shares
		from holding h
		left join company c on h.company_id = c.id
		where h.trader_id = %v`, hs.uid)
}