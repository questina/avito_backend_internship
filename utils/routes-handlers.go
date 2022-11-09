package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/questina/avito_backend_internship/db"
	"github.com/questina/avito_backend_internship/schemes"
	"net/http"
	"os"
	"strconv"
)

// AddMoney godoc
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
func AddMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var up schemes.UpgradeBalance
	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var id = UpdateBalance(up.Amount, up.Id)
	if id == -1 {
		http.Error(rw, "Could not write to database", http.StatusInternalServerError)
		return
	}
	var userBalance = GetUserBalance(id, true)
	if userBalance == -1 {
		http.Error(rw, "Could not read user id from database", http.StatusInternalServerError)
		return
	}
	AddEvent("ADD", up.Amount, -1, -1, id)
	var resp = schemes.AddMoneyReturn{Id: id, Balance: userBalance, Status: "OK"}
	json.NewEncoder(rw).Encode(resp)
}

// ReserveMoney godoc
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
func ReserveMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var ordReserve schemes.OrderReserve
	err := json.NewDecoder(r.Body).Decode(&ordReserve)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var userBalance = GetUserBalance(ordReserve.UserId, false)
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

	tx, err := db.Db.Begin()

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
	AddEvent("RESERVE", ordReserve.Cost, ordReserve.ServiceId, ordReserve.OrderId, ordReserve.UserId)
	json.NewEncoder(rw).Encode(map[string]string{"Status": "OK"})
}

// TakeMoney godoc
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
func TakeMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var ordReserve schemes.OrderReserve
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
	var userBalance = GetUserBalance(ordReserve.UserId, true)
	if userBalance == -1 {
		http.Error(rw, "Could not read user id from database", http.StatusInternalServerError)
		return
	}
	tx, err := db.Db.Begin()
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
	AddEvent("TAKE", ordReserve.Cost, ordReserve.ServiceId, ordReserve.OrderId, ordReserve.UserId)
	json.NewEncoder(rw).Encode(map[string]string{"Status": "OK"})
}

// GetBalance getBalance godoc
// @Summary     Get user balance
// @Description Get user balance from user account
// @Accept      json
// @Produce     json
// @Param       id     body     int   true "User ID"
// @Success     200    {object} User
// @Failure 400
// @Failure 500
// @Router      /get_balance [post]
func GetBalance(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var user schemes.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	user.Balance = GetUserBalance(user.Id, false)
	if user.Balance == -1 {
		http.Error(rw, "Could not read user from database", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(user)
}

// FreeMoney freeMoney godoc
// @Summary     Free reserved money
// @Description Free reserved money from user account
// @Accept      json
// @Produce     json
// @Param       UserId     body     int   true "User ID"
// @Param       OrderId     body int true  "Order ID"
// @Param		ServiceId	body int true "Service ID"
// @Param		Cost body number true "Order cost"
// @Success     200    {object} StatusMessage
// @Failure 400
// @Failure 500
// @Router      /free_money [post]
func FreeMoney(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var ordReserve schemes.OrderReserve
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

	tx, err := db.Db.Begin()
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
	AddEvent("FREE", ordReserve.Cost, ordReserve.ServiceId, ordReserve.OrderId, ordReserve.UserId)
	json.NewEncoder(rw).Encode(map[string]string{"Status": "OK"})
}

// GenReport genReport godoc
// @Summary     Generate report
// @Description Generate report on monthly income from services
// @Accept      json
// @Produce     octet-stream
// @Param       Year     body     int   true "Search Year"
// @Param       Month     body int true  "Search Month"
// @Success     200    {object} StatusMessage
// @Failure 400
// @Failure 500
// @Router      /generate_report [post]
func GenReport(rw http.ResponseWriter, r *http.Request) {
	var reportInfo schemes.ReportInput
	err := json.NewDecoder(r.Body).Decode(&reportInfo)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fmt.Sprintf("./reports/report_%d_%d.csv", reportInfo.Month, reportInfo.Year)))
	rw.Header().Set("Content-Type", "application/octet-stream")
	var (
		service_income schemes.ReportData
		data           []schemes.ReportData
	)
	rows, err := db.Db.Query("SELECT service_id, SUM(amount) FROM moneyflow "+
		"WHERE event_type=\"TAKE\" AND YEAR(datetime)=? AND MONTH(datetime)=? GROUP BY service_id",
		reportInfo.Year, reportInfo.Month)
	if err != nil {
		http.Error(rw, "Could not get info from database", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		rows.Scan(&service_income.Service_id, &service_income.Income)
		data = append(data, service_income)
	}
	defer rows.Close()

	f, err := os.Create(fmt.Sprintf("./reports/report_%d_%d.csv", reportInfo.Month, reportInfo.Year))
	if err != nil {
		http.Error(rw, "Failed to open file", http.StatusInternalServerError)
		return
	}
	w := csv.NewWriter(f)
	for _, service_line := range data {
		row := []string{strconv.Itoa(service_line.Service_id), strconv.FormatFloat(float64(service_line.Income), 'f', -1, 32)}
		if err := w.Write(row); err != nil {
			http.Error(rw, "Failed to write into file", http.StatusInternalServerError)
			return
		}
	}
	w.Flush()
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.ServeFile(rw, r, fmt.Sprintf("./reports/report_%d_%d.csv", reportInfo.Month, reportInfo.Year))
	return
}

// BalanceInfo godoc
// @Summary     Give balance info
// @Description Give balance info on user expenses
// @Accept      json
// @Produce     json
// @Param       UserId     body     int   true "User ID"
// @Param   sort  query     string     false  "Direction of Sorting"       Enums(asc, desc)
// @Param   limit  query     int     false  "Search Limit"
// @Param   offset  query     int    false  "Search offset"
// @Success     200    {object} StatusMessage
// @Failure 400
// @Failure 500
// @Router      /balance_info [post]
func BalanceInfo(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var userId schemes.UserId
	err := json.NewDecoder(r.Body).Decode(&userId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(userId.Id)
	var (
		limit  = "10"
		offset = "0"
		sort   = "asc"
	)
	q := r.URL.Query()
	if len(q["limit"]) != 0 {
		limit = q["limit"][0]
	}
	if len(q["offset"]) != 0 {
		offset = q["offset"][0]
	}
	if len(q["sort"]) != 0 {
		sort = q["sort"][0]
	}
	var (
		event  schemes.BalanceInfo
		events []schemes.BalanceInfo
	)
	var query string
	if sort == "asc" {
		query = "SELECT datetime, amount, event_type, service_id, order_id FROM moneyflow " +
			"WHERE user_id=? ORDER BY datetime ASC LIMIT ? OFFSET ?"
	} else {
		query = "SELECT datetime, amount, event_type, service_id, order_id FROM moneyflow " +
			"WHERE user_id=? ORDER BY datetime DESC LIMIT ? OFFSET ?"
	}
	rows, err := db.Db.Query(query, userId.Id, limit, offset)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Could not get info from database", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		rows.Scan(&event.Timestamp, &event.Amount, &event.EventType, &event.ServiceId, &event.OrderId)
		events = append(events, event)
	}
	defer rows.Close()
	json.NewEncoder(rw).Encode(events)
}
