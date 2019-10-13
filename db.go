package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `CREATE TABLE IF NOT EXISTS userinfo (
	uuid CHAR(32) PRIMARY KEY NOT NULL,
	username VARCHAR(64) NOT NULL,
	nickname VARCHAR(64) NULL,
	created DATE NULL
)`

func openDatabase() (db *sql.DB, err error) {
	defer func() {
		if err != nil {
			log.Printf("OpenDatabase error: %s", err)
			err = fmt.Errorf("OpenDatabase: %s", err)
		}
	}()

	db, err = sql.Open("sqlite3", "./ayame.db")
	stmt, err := db.Prepare(schema)
	_, err = stmt.Exec()

	return
}

func readTable(db *sql.DB) (err error) {
	return
}

func readRecord(db *sql.DB) (err error) {
	return
}

func addRecord(db *sql.DB) (err error) {
	return
}

func deleteRecord(db *sql.DB) (err error) {
	return
}

func updateRecord(db *sql.DB) (err error) {
	return
}
