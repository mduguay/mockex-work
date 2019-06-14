package main

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

func (s *Storage) connect() {
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

func (s *Storage) disconnect() {
	fmt.Println("Disconnecting")
	s.db.Close()
}

var (
	id   int
	name string
)

func (s *Storage) users() {
	rows, err := s.db.Query("select id, email from users where id = $1", 1)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Println(err)
		}
		log.Println(id, name)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
}

func (s *Storage) readUser(id int, result chan string) {
	fmt.Println("Preparing")
	stmt, err := s.db.Prepare("select email from users where id = $1")
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
