package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type db struct {
	host     string
	port     int
	user     string
	password string
	name     string
}

func main() {

	dbInfo := db{"localhost", 5432, "postgres", "postgres", "cdb"}
	psql := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbInfo.host, dbInfo.port, dbInfo.user, dbInfo.password, dbInfo.name)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
}
