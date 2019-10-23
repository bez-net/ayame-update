package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `CREATE TABLE IF NOT EXISTS userinfo (
	uuid CHAR(32) PRIMARY KEY NOT NULL,
	name VARCHAR(64) NOT NULL,
	nick VARCHAR(64) NULL,
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

func readTable(db *sql.DB, user *User) (err error) {
	rows, err := db.Query("SELECT * FROM userinfo")
	for rows.Next() {
		err = rows.Scan(&user.uuid, &user.name, &user.nick, &user.created)
		if err != nil {
			break
		}
	}
	return
}

func readRecord(db *sql.DB, user *User) (err error) {
	stmt, err := db.Prepare("SELECT DISTINCT * FROM userinfo where uuid=?")
	_, err = stmt.Exec(user.uuid)
	return
}

func addRecord(db *sql.DB, user *User) (err error) {
	stmt, err := db.Prepare("INSERT INTO userinfo(uuid, name, nick, created) values(?,?,?,?)")
	_, err = stmt.Exec(user.uuid, user.name, user.nick, user.created)
	return
}

func deleteRecord(db *sql.DB, user *User) (err error) {
	stmt, err := db.Prepare("DELETE from userinfo where uuid=?")
	_, err = stmt.Exec(user.uuid)
	return
}

func updateRecord(db *sql.DB, user *User) (err error) {
	stmt, err := db.Prepare("UPDATE userinfo set name=? where uuid=?")
	_, err = stmt.Exec("stoney", user.uuid)
	return
}
