package main

import (
	"encoding/json"
	"net/http"
)

// addMoney godoc
// @Summary     Add money to user account
// @Description Create new user id if user does not exists and add money to his account
// @Accept      json
// @Produce     json
// @Param       id     body     int   false "User ID"
// @Param       amount body     number true  "Amount of money"
// @Success     200    {object} AddMoneyReturn
// @Failure 400
// @Failure 500
// @Router      /add_money [post]
func addMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var up UpgradeBalance
	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var id = updateBalance(up.Amount, up.Id)
	if id == -1 {
		http.Error(rw, "Could not write to database", http.StatusInternalServerError)
		return
	}
	var userBalance = getUserBalance(id)
	if userBalance == -1 {
		http.Error(rw, "Could not read user id from database", http.StatusInternalServerError)
		return
	}
	var resp = AddMoneyReturn{Id: id, Balance: userBalance, Status: "OK"}
	json.NewEncoder(rw).Encode(resp)
}

// reserveMoney godoc
// @Summary     Reserve money on user account
// @Description Get request from service and reserve money on user account
// @Accept      json
// @Produce     json
// @Param       UserId     body     int   true "User ID"
// @Param       OrderId     body int true  "Order ID"
// @Param		ServiceId	body int true "Service ID"
// @Param		Cost body number true "Order cost"
// @Success     200    {object} StatusMessage
// @Failure 400
// @Failure 500
// @Router      /reserve_money [post]
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
		http.Error(rw, "Could not read user id from database", http.StatusInternalServerError)
		return
	}

	if userBalance < ordReserve.Cost {
		http.Error(rw, "User does not have enough money", http.StatusBadRequest)
		return
	}

	var order_exists = CheckOrderId(ordReserve.OrderId)
	if order_exists {
		http.Error(rw, "Order already exists", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()

	stmt, err := tx.Prepare("INSERT into orders SET order_id=?,service_id=?,user_id=?,cost=?")
	if err != nil {
		tx.Rollback()
		http.Error(rw, "Could not write new order to database", http.StatusInternalServerError)
		return
	}
	_, queryError := stmt.Exec(ordReserve.OrderId, ordReserve.ServiceId, ordReserve.UserId, ordReserve.Cost)
	if queryError != nil {
		tx.Rollback()
		http.Error(rw, "Could not write new order to database", http.StatusInternalServerError)
		return
	}

	success := UpdateReservedBalance(ordReserve.Cost, ordReserve.UserId, tx)
	if !success {
		tx.Rollback()
		http.Error(rw, "Could not update reserved balance to database", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(rw, "Could not commit changes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(map[string]string{"Status": "Ok"})
}

// takeMoney godoc
// @Summary     Take reserved money
// @Description Take reserved money from user account
// @Accept      json
// @Produce     json
// @Param       UserId     body     int   true "User ID"
// @Param       OrderId     body int true  "Order ID"
// @Param		ServiceId	body int true "Service ID"
// @Param		Cost body number true "Order cost"
// @Success     200    {object} StatusMessage
// @Failure 400
// @Failure 500
// @Router      /take_money [post]
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
		http.Error(rw, "Order does not exists", http.StatusBadRequest)
		return
	}
	var userBalance = getUserBalance(ordReserve.UserId)
	if userBalance == -1 {
		http.Error(rw, "Could not read user id from database", http.StatusInternalServerError)
		return
	}
	tx, err := db.Begin()
	_, err = tx.Exec("DELETE FROM orders WHERE order_id=?", ordReserve.OrderId)
	if err != nil {
		tx.Rollback()
		http.Error(rw, "Could not delete order from database", http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec("UPDATE user_balances SET balance=? WHERE id=?", userBalance-ordReserve.Cost, ordReserve.UserId)
	if err != nil {
		tx.Rollback()
		http.Error(rw, "Could not update user balance from database", http.StatusInternalServerError)
		return
	}

	success := UpdateReservedBalance((-1)*ordReserve.Cost, ordReserve.UserId, tx)
	if !success {
		tx.Rollback()
		http.Error(rw, "Could not update reserved balance to database", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(rw, "Could not commit changes", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(map[string]string{"Status": "Ok"})
}

// getBalance godoc
// @Summary     Get user balance
// @Description Get user balance from user account
// @Accept      json
// @Produce     json
// @Param       id     body     int   true "User ID"
// @Success     200    {object} User
// @Failure 400
// @Failure 500
// @Router      /get_balance [post]
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
		http.Error(rw, "Could not read user from database", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(user)
}

func freeMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var ordReserve OrderReserve
	err := json.NewDecoder(r.Body).Decode(&ordReserve)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var orderExists = CheckOrderId(ordReserve.OrderId)
	if !orderExists {
		http.Error(rw, "Order does not exists", http.StatusBadRequest)
		return
	}
	var userBalance = getUserBalance(ordReserve.UserId)
	if userBalance == -1 {
		http.Error(rw, "Could not read user id from database", http.StatusInternalServerError)
		return
	}
	tx, err := db.Begin()
	_, err = tx.Exec("DELETE FROM orders WHERE order_id=?", ordReserve.OrderId)
	if err != nil {
		tx.Rollback()
		http.Error(rw, "Could not delete order from database", http.StatusInternalServerError)
		return
	}

	success := UpdateReservedBalance((-1)*ordReserve.Cost, ordReserve.UserId, tx)
	if !success {
		tx.Rollback()
		http.Error(rw, "Could not update reserved balance to database", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(rw, "Could not commit changes", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(map[string]string{"Status": "Ok"})
}
