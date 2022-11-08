package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func getUsers(rw http.ResponseWriter, r *http.Request) {
	var (
		user  User
		users []User
	)
	rows, err := db.Query("SELECT * FROM user_balances")
	if err != nil {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not read from database. Try later"})
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Balance)
		users = append(users, user)
	}
	defer rows.Close()
	json.NewEncoder(rw).Encode(users)
}

func addMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var up UpgradeBalance
	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var success = updateBalance(up.Amount, up.Id)
	if !success {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not write to database. Try later"})
		return
	}
	json.NewEncoder(rw).Encode(map[string]string{"Status": "OK"})
}

func reserveMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var ordReserve OrderReserve
	err := json.NewDecoder(r.Body).Decode(&ordReserve)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO check if order id already exists
	stmt, err := db.Prepare("INSERT into orders SET order_id=?,service_id=?,user_id=?,cost=?")
	if err != nil {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not write to database, please try later"})
		fmt.Println(err)
		return
	}
	_, queryError := stmt.Exec(ordReserve.OrderId, ordReserve.ServiceId, ordReserve.UserId, ordReserve.Cost)
	if queryError != nil {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not write to database, please try later"})
		fmt.Println(err)
		return
	}
	json.NewEncoder(rw).Encode(map[string]string{"Status": "Ok"})
}

func takeMoney(rw http.ResponseWriter, r *http.Request) {
	// TODO check if order exists in database
	rw.Header().Set("Content-Type", "application/json")
	var ordReserve OrderReserve
	err := json.NewDecoder(r.Body).Decode(&ordReserve)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var userBalance int
	if err := db.QueryRow("SELECT balance from user_balances where id = ?",
		ordReserve.UserId).Scan(&userBalance); err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not write to database, please try later"})
			fmt.Println(err)
			return
		}
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not write to database, please try later"})
		fmt.Println(err)
		return
	}

	tx, err := db.Begin()
	_, err = tx.Exec("DELETE FROM orders WHERE order_id=?", ordReserve.OrderId)
	if err != nil {
		tx.Rollback()
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not write to database, please try later"})
		fmt.Println(err)
		return
	}
	_, err = tx.Exec("UPDATE user_balances SET balance=? WHERE id=?", userBalance-ordReserve.Cost, ordReserve.UserId)
	if err != nil {
		tx.Rollback()
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not write to database, please try later"})
		fmt.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not commit changes, please try later"})
		fmt.Println(err)
		return
	}
	json.NewEncoder(rw).Encode(map[string]string{"Status": "Ok"})
}

func getBalance(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(user.Id)
	var userBalance int
	if err := db.QueryRow("SELECT balance from user_balances where id = ?",
		user.Id).Scan(&userBalance); err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not read from database, please try later"})
			fmt.Println(err)
			return
		}
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not read from database, please try later"})
		fmt.Println(err)
		return
	}
	user.Balance = userBalance
	json.NewEncoder(rw).Encode(user)
}

func updateBalance(newBalance int, userId int) bool {
	var userBalance int
	if err := db.QueryRow("SELECT balance from user_balances where id = ?",
		userId).Scan(&userBalance); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
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
