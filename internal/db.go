package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	// Blank import of the postgres sql driver
	_ "github.com/lib/pq"
)

const (
	connstring = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
)

// Storage is the object that connects and interacts with the database
type Storage struct {
	db *sql.DB
}

// Connect initializes storage and establishes a connection to the db
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

// Disconnect will break the connection with the db
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
		ID: id,
	}
	err = stmt.QueryRow(id).Scan(&t.Email)
	check(err)
	return t
}

func (s *Storage) readCash(tid int) Cash {
	stmt, err := s.db.Prepare("select amount from cash where trader_id = $1")
	check(err)
	defer stmt.Close()
	c := Cash{
		Tid: tid,
	}
	err = stmt.QueryRow(tid).Scan(&c.Amount)
	check(err)
	return c
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

func (s *Storage) createTrade(t *Trade) (*TradeResult, error) {
	tx, err := s.db.Begin()
	if checktx(err, tx) {
		return nil, err
	}

	var shares int
	{
		stmt, err := tx.Prepare("SELECT shares FROM holding WHERE company_id = $1 AND trader_id = $2")
		if checktx(err, tx) {
			return nil, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(t.Cid, t.Tid).Scan(&shares)
		if checktx(err, tx) {
			return nil, err
		}
	}

	var amount float64
	{
		stmt, err := tx.Prepare("SELECT amount FROM cash WHERE trader_id = $1")
		if checktx(err, tx) {
			return nil, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(t.Tid).Scan(&amount)
		if checktx(err, tx) {
			return nil, err
		}
	}

	amount -= float64(t.Shares) * t.Price
	shares += t.Shares
	if amount < 0 {
		tx.Rollback()
		return nil, errors.New("Not enough cash to buy")
	}
	if shares < 0 {
		tx.Rollback()
		return nil, errors.New("Not enough shares to sell")
	}

	{
		stmt, err := tx.Prepare("INSERT INTO trade(trader_id, company_id, shares, price) VALUES($1, $2, $3, $4)")
		if checktx(err, tx) {
			return nil, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(t.Tid, t.Cid, t.Shares, t.Price)
		if checktx(err, tx) {
			return nil, err
		}
	}
	{
		stmt, err := tx.Prepare("UPDATE holding SET shares = $1 WHERE trader_id = $2 AND company_id = $3")
		if checktx(err, tx) {
			return nil, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(shares, t.Tid, t.Cid)
		if checktx(err, tx) {
			return nil, err
		}
	}
	{
		stmt, err := tx.Prepare("UPDATE cash SET amount = $1 WHERE trader_id = $2")
		if checktx(err, tx) {
			return nil, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(amount, t.Tid)
		if checktx(err, tx) {
			return nil, err
		}
	}

	err = tx.Commit()
	if checktx(err, tx) {
		return nil, err
	}

	tradeResult := &TradeResult{
		Shares: shares,
		Amount: amount,
	}
	return tradeResult, nil
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

func (s Storage) updateSettings(cid int, settings *Settings) {
	ustmt, err := s.db.Prepare("update stock set price = $1, vol = $2, minchange = $3, maxchange = $4 where company_id = $5")
	check(err)
	r, err := ustmt.Exec(settings.Price, settings.Vol, settings.Minchange, settings.Maxchange, cid)
	check(err)
	fmt.Println("updateSettings Result:", r)
	return
}

func checktx(err error, tx *sql.Tx) bool {
	if check(err) {
		tx.Rollback()
		return true
	}
	return false
}
