package data

import (
	"database/sql"
	"fmt"
)

// Trader is a user that holds stock and makes trades
type Trader struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// TraderScanner is responsible for fetching a single trader given an id
type TraderScanner struct {
	id int
}

// Query is the db query to be executed
func (ts *TraderScanner) Query() string {
	return fmt.Sprintf("select email from trader where id = %v", ts.id)
}

// ScanRow handles the result given from the db query above
func (ts *TraderScanner) ScanRow(stmt *sql.Stmt) (interface{}, error) {
	t := new(Trader)
	err := stmt.QueryRow(ts.id).Scan(t.Email)
	return t, err
}
