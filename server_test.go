package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var route mux.Router

func TestMain(m *testing.M) {
	LoadEnv()

	fmt.Println("Server will start at http://localhost:8000/")

	connectDatabase()

	route := mux.NewRouter()

	addApproutes(route)
	code := m.Run()
	os.Exit(code)
}

func TestGetUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUsers)

	handler.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusOK, rr.Code)
}

func TestAddUser(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	if body := rr.Body.String(); !strings.Contains(body, "\"Status\":\"OK\"") {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	var user User
	json.NewDecoder(rr.Body).Decode(&user)
	fmt.Println(user)
	jsonBody = []byte(fmt.Sprintf(`{"id": %d, "amount": 300}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/add_money", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(addMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	if body != fmt.Sprintf("{\"Balance\": \"%f\", \"Id\": \"%d\", \"Status\": \"OK\"}", float32(400), user.Id) {
		clearTable(user.Id, 0)
		t.Errorf("Returned status not OK. Got %s", body)
	}
	clearTable(user.Id, 0)
}

func TestReserveMoney(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	if !strings.Contains(body, "'Status': 'OK'") {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	var user User
	json.NewDecoder(rr.Body).Decode(&user)

	jsonBody = []byte(fmt.Sprintf(`{
		"UserId": %d,
		"OrderId": 1,
		"ServiceId": 1,
		"Cost": 10
	}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/reserve_money", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(reserveMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	fmt.Println(rr.Body.String())
	body = rr.Body.String()
	if body != "{\"Status\":\"OK\"}" {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	clearTable(user.Id, 1)
}

func TestTakeMoney(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	if strings.Contains(body, "'Status': 'OK'") {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	var user User
	json.NewDecoder(rr.Body).Decode(&user)

	jsonBody = []byte(fmt.Sprintf(`{
		"UserId": %d,
		"OrderId": 1,
		"ServiceId": 1,
		"Cost": 10
	}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/reserve_money", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(reserveMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	if body := rr.Body.String(); strings.Contains(body, "'Status': 'OK'") {
		t.Errorf("Returned status not OK. Got %s", body)
	}

	jsonBody = []byte(fmt.Sprintf(`{
		"UserId": %d,
		"OrderId": 1,
		"ServiceId": 1,
		"Cost": 10
	}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/take_money", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(takeMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	if body := rr.Body.String(); strings.Contains(body, "'Status': 'OK'") {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	clearTable(user.Id, 1)
}

func TestGetBalance(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	if strings.Contains(body, "'Status': 'OK'") {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	var user User
	jsonBody = []byte(fmt.Sprintf(`{"id": %d}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/get_balance", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(getBalance)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body = rr.Body.String()
	if body == fmt.Sprintf("{\"Id\": %d, \"Balance\": %f}", user.Id, user.Balance) {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	clearTable(user.Id, 0)
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTable(user_id int, order_id int) {
	if order_id > 0 {
		stmt, _ := db.Prepare("DELETE FROM orders WHERE order_id=?")
		_, err := stmt.Exec(order_id)
		if err != nil {
			fmt.Printf("Could not delete from orders with order_id=%d\n", order_id)
		}
	}
	if user_id > 0 {
		stmt, _ := db.Prepare("DELETE FROM user_balances WHERE user_id=?")
		_, err := stmt.Exec(order_id)
		if err != nil {
			fmt.Printf("Could not delete from user_balances with user_id=%d\n", user_id)
		}
	}
}
