package data

import (
	"database/sql"
	"fmt"
)

// Cash is the amount of money on hand that can be used to buy stocks
type Cash struct {
	Tid    int     `json:"tid"`
	Amount float64 `json:"amount"`
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
func (cs *CashScanner) ScanRow(stmt *sql.Stmt) (interface{}, error) {
	c := new(Cash)
	err := stmt.QueryRow(cs.id).Scan(c.Amount)
	return c, err
}
