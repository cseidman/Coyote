package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func (v *VM) ExecSQL(sql string) {

}

func OpenDb() *sql.DB {
	db, _ := sql.Open("sqlite3", ":memory:")
	return db
}

func CreateDb() {
	//OpenDb()
}