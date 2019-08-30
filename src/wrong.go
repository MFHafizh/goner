package src

import (
	"database/sql"
	"fmt"
	"log"
)

type textRange struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine"`
	StartColumn int `json:"startColumn,omitempty"`
	EtartColumn int `json:"endColumn,omitempty"`
}

func buildSqlUsingStrFormat(email string) *sql.Row {
	sqlStatement := fmt.Sprintf("SELECT * FROM users WHERE email='%s';", email)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		"localhost", "5432", "user", "dbname")
	log.Println("connString", psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	row := db.QueryRow(sqlStatement)
	return row
}

func buildSqlUsingStrConcate(email string) *sql.Row {
	sqlStatement := "SELECT * FROM users WHERE email='" + email + "';"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		"localhost", "5432", "user", "dbname")
	log.Println("connString", psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	row := db.QueryRow(sqlStatement)
	return row
}
