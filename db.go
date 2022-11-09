package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB
var err error

func connectDatabase() {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")
	fmt.Println(DBNAME)
	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s", USER, PASS, HOST, DBNAME)
	fmt.Println(URL)
	db, err = sql.Open("mysql", URL)
	if err != nil {
		fmt.Println("Can't conect to database")
	} else {
		fmt.Println("Database connected.")
	}
}
