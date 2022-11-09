package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var Db *sql.DB
var err error

func ConnectDatabase() {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")
	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s", USER, PASS, HOST, DBNAME)
	Db, err = sql.Open("mysql", URL)
	if err != nil {
		fmt.Println("Can't conect to database")
	} else {
		fmt.Println("Database connected.")
	}
}
