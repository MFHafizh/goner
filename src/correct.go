package src

import (
	"database/sql"
	"fmt"
	"log"
)

func buildSQL(email string) *sql.Row {
	sqlStatement := `SELECT id, name, email FROM public."member" WHERE id=$1;`
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
	row := db.QueryRow(sqlStatement, "mail")
	return row
}
