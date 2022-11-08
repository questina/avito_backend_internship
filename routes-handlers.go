package main

import (
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

	var userBalance = getUserBalance(ordReserve.UserId)
	if userBalance == -1 {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not read user id from database, please try later"})
		fmt.Println(err)
		return
	}

	if userBalance < ordReserve.Cost {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "User does not have enough money"})
		fmt.Println(err)
		return
	}

	var order_exists = CheckOrderId(ordReserve.OrderId)
	if order_exists {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Order already exists"})
		fmt.Println(err)
		return
	}
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
	rw.Header().Set("Content-Type", "application/json")
	var ordReserve OrderReserve
	err := json.NewDecoder(r.Body).Decode(&ordReserve)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var orderExists = CheckOrderId(ordReserve.OrderId)
	if !orderExists {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Order does not exists"})
		fmt.Println(err)
		return
	}
	var userBalance = getUserBalance(ordReserve.UserId)
	if userBalance == -1 {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not read user id from database, please try later"})
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
	user.Balance = getUserBalance(user.Id)
	if user.Balance == -1 {
		json.NewEncoder(rw).Encode(map[string]string{"Status": "Could not read user id from database, please try later"})
		fmt.Println(err)
		return
	}
	json.NewEncoder(rw).Encode(user)
}
