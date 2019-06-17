package main

import "database/sql"

type Scanner interface {
	ScanRow(rows *sql.Rows) interface{}
}

type CompanyScanner struct{}

func (cs *CompanyScanner) ScanRow(rows *sql.Rows) interface{} {
	c := new(Company)
	err := rows.Scan(&c.Id, &c.Symbol)
	check(err)
	return c
}
