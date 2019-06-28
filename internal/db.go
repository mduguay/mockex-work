package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	connstring = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Connect() {
	env := getEnv()
	log.Println("DB: Connecting")
	conString := fmt.Sprintf(connstring, env["host"], env["port"], env["user"], env["pass"], env["name"])
	db, err := sql.Open("postgres", conString)
	check(err)
	err = db.Ping()
	check(err)
	log.Println("DB: Connected")
	s.db = db
}

func (s *Storage) Disconnect() {
	log.Println("DB: Disconnecting")
	s.db.Close()
}

func getEnv() map[string]string {
	env := make(map[string]string)
	env["host"] = os.Getenv("DB_HOST")
	env["port"] = os.Getenv("DB_PORT")
	env["user"] = os.Getenv("DB_USER")
	env["pass"] = os.Getenv("DB_PASS")
	env["name"] = os.Getenv("DB_NAME")
	return env
}

func (s *Storage) readTrader(id int) Trader {
	stmt, err := s.db.Prepare("select email from trader where id = $1")
	check(err)
	defer stmt.Close()
	t := Trader{
		Id: id,
	}
	err = stmt.QueryRow(id).Scan(&t.Email)
	check(err)
	return t
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
	var shares int
	tx, err := s.db.Begin()
	if checktx(err, tx) {
		return
	}
	{
		stmt, err := tx.Prepare("SELECT shares FROM holding WHERE company_id = $1 AND trader_id = $2")
		if checktx(err, tx) {
			return
		}
		defer stmt.Close()

		err = stmt.QueryRow(t.Cid, t.Tid).Scan(&shares)
		if checktx(err, tx) {
			return
		}
	}
	if shares+t.Shares < 0 {
		tx.Rollback()
		return
	}
	{
		stmt, err := tx.Prepare("INSERT INTO trade(trader_id, company_id, direction, shares, price) VALUES($1, $2, $3, $4, $5)")
		if checktx(err, tx) {
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(t.Tid, t.Cid, t.Direction, t.Shares, t.Price)
		if checktx(err, tx) {
			return
		}
	}
	shares += t.Shares
	{
		stmt, err := tx.Prepare("UPDATE holding SET shares = $1 WHERE trader_id = $2 AND company_id = $3")
		if checktx(err, tx) {
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(shares, t.Tid, t.Cid)
		if checktx(err, tx) {
			return
		}
	}
	err = tx.Commit()
	if checktx(err, tx) {
		return
	}
}

func (s *Storage) createQuote(q *Quote) {
	var id int
	tx, err := s.db.Begin()
	if checktx(err, tx) {
		return
	}
	{
		stmt, err := tx.Prepare("SELECT id FROM company WHERE symbol = $1")
		if checktx(err, tx) {
			return
		}
		defer stmt.Close()
		err = stmt.QueryRow(q.Symbol).Scan(&id)
		if checktx(err, tx) {
			return
		}
	}
	{
		istmt, err := tx.Prepare("INSERT INTO quote (company_id, price, stamp) VALUES ($1, $2, $3)")
		if checktx(err, tx) {
			return
		}
		defer istmt.Close()
		_, err = istmt.Exec(id, q.Price, time.Now())
		if checktx(err, tx) {
			return
		}
	}
	err = tx.Commit()
	if checktx(err, tx) {
		return
	}
}

func checktx(err error, tx *sql.Tx) bool {
	if check(err) {
		tx.Rollback()
		return true
	}
	return false
}
