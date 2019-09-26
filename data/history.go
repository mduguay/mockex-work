package data

import (
	"database/sql"
	"fmt"
	"time"
)

// HistoryScanner is responsible for fetching the history of a given stock
type HistoryScanner struct {
	cid int
}

// HistPoint is one point in the history of a stock's price
type HistPoint struct {
	Price float64   `json:"price"`
	Stamp time.Time `json:"stamp"`
}

// Query is the db query to be executed
func (hs *HistoryScanner) Query() string {
	t := time.Now()
	location, _ := time.LoadLocation("America/New_York")
	opening := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, location)
	fmt.Println(opening)
	return fmt.Sprintf("select price, stamp from quote where company_id = %v and stamp > now()::date order by stamp desc", hs.cid)
}

// ScanRow reads the results from storage and creates a HistPoint
func (hs *HistoryScanner) ScanRow(rows *sql.Rows) (interface{}, error) {
	hp := new(HistPoint)
	err := rows.Scan(&hp.Price, &hp.Stamp)
	return hp, err
}
