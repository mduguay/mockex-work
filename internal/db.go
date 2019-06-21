package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	connstring = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	host       = "localhost"
	port       = 5432
	user       = "postgres"
	password   = "postgres64"
	dbname     = "mockex"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Connect() {
	fmt.Println("Connecting")
	conString := fmt.Sprintf(connstring, host, port, user, password, dbname)
	db, err := sql.Open("postgres", conString)
	check(err)
	err = db.Ping()
	check(err)
	s.db = db
}

func (s *Storage) Disconnect() {
	fmt.Println("Disconnecting")
	s.db.Close()
}

func (s *Storage) readTrader(id int, result chan string) {
	fmt.Println("Preparing")
	stmt, err := s.db.Prepare("select email from trader where id = $1")
	check(err)
	defer stmt.Close()
	fmt.Println("Querying")
	var name string
	err = stmt.QueryRow(id).Scan(&name)
	check(err)
	fmt.Println("Publishing")
	fmt.Println(name)
	result <- name
}

func (s *Storage) readMultiple(scanner Scanner) (items []interface{}) {
	stmt, err := s.db.Prepare(scanner.Query())
	check(err)
	defer stmt.Close()
	rows, err := stmt.Query()
	check(err)
	defer rows.Close()
	for rows.Next() {
		item := scanner.ScanRow(rows)
		items = append(items, item)
	}
	err = rows.Err()
	check(err)
	return
}

func (s *Storage) createTrade(t Trade) {
	// Create trade
	stmt, err := s.db.Prepare("INSERT INTO trade(trader_id, company_id, direction, shares, price) VALUES($1, $2, $3, $4, $5)")
	check(err)
	res, err := stmt.Exec(t.Tid, t.Sid, t.Direction, t.Shares, t.Price)
	fmt.Printf("%+v", res)
	// Update holding
}
