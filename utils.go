package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("unable to load .env file")
	}
}

func getUserBalance(userId int) float32 {
	var userBalance float32
	if err := db.QueryRow("SELECT balance from user_balances where id = ?",
		userId).Scan(&userBalance); err != nil {
		if err == sql.ErrNoRows {
			return -1
		}
		return -1
	}
	return userBalance
}

func updateBalance(newBalance float32, userId int) bool {
	var userBalance = getUserBalance(userId)
	if userBalance == -1 {
		stmt, err := db.Prepare("INSERT into user_balances SET balance=?")
		if err != nil {
			fmt.Println(err)
			return false
		}
		_, queryError := stmt.Exec(newBalance)
		if queryError != nil {
			fmt.Println(queryError)
			return false
		}
	} else {
		stmt, err := db.Prepare("UPDATE user_balances SET balance=? WHERE id=?")
		if err != nil {
			fmt.Println(err)
			return false
		}
		_, queryError := stmt.Exec(userBalance+newBalance, userId)
		if queryError != nil {
			fmt.Println(queryError)
			return false
		}
	}
	return true
}

func CheckOrderId(orderId int) bool {
	var servId int
	if err := db.QueryRow("SELECT service_id from orders where order_id=?",
		orderId).Scan(&servId); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	return true
}
