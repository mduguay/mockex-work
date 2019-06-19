package internal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres64"
	dbname   = "mockex"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Connect() {
	fmt.Println("Connecting")
	conString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", conString)
	if err != nil {
		log.Println(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	s.db = db
}

func (s *Storage) Disconnect() {
	fmt.Println("Disconnecting")
	s.db.Close()
}

var (
	id   int
	name string
)

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
	if err = rows.Err(); err != nil {
		check(err)
	}
	return
}
