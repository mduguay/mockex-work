package data

import (
	"database/sql"
	"fmt"
)

// Holding represents a traders assets in a particular stock
type Holding struct {
	Tid    int    `json:"tid"`
	Cid    int    `json:"cid"`
	Symbol string `json:"symbol"`
	Shares int    `json:"shares"`
}

// HoldingScanner is responsible for fetching the holdings for a given trader
type HoldingScanner struct {
	UID string
}

// Query is the db query to be executed
func (hs *HoldingScanner) Query() string {
	return fmt.Sprintf(`
		select h.trader_id, c.id, c.symbol, h.shares
		from holding h
		left join company c on h.company_id = c.id
		where h.trader_id = %v`, hs.UID)
}

// ScanRow reads from the db results and creates a Holding
func (hs *HoldingScanner) ScanRow(rows *sql.Rows) (interface{}, error) {
	h := new(Holding)
	err := rows.Scan(&h.Tid, &h.Cid, &h.Symbol, &h.Shares)
	return h, err
}
