package main

import (
	"database/sql"
	"fmt"
)

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
		return false
	}
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
	return true
}

func CheckOrderId(orderId int) bool {
	var order OrderReserve
	if err := db.QueryRow("SELECT * from orders where order_id = ?",
		orderId).Scan(&order); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	return true
}
