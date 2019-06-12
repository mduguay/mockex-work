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

func (s *Storage) singleUser() {
	//statement := `SELECT id, email FROM users WHERE id=$1`
}
