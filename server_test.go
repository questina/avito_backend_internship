package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/questina/avito_backend_internship/db"
	"github.com/questina/avito_backend_internship/schemes"
	"github.com/questina/avito_backend_internship/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	utils.Start()
	code := m.Run()
	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.AddMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	if body := rr.Body.String(); !strings.Contains(body, "\"Status\":\"OK\"") {
		t.Errorf("Returned status not OK. Got %s", body)
	}
	var user schemes.User
	json.NewDecoder(rr.Body).Decode(&user)
	jsonBody = []byte(fmt.Sprintf(`{"id": %d, "amount": 300}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/add_money", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(utils.AddMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	var mapBody map[string]interface{}
	json.Unmarshal([]byte(body), &mapBody)
	if money, ok := mapBody["Balance"]; ok {
		if money != float64(400) {
			clearTable(user.Id, 0)
			t.Errorf("Incorrect balance added. Expected %f, got %f", float32(400), money)
		}
	} else {
		clearTable(user.Id, 0)
		t.Errorf("Incorrect json output in add_money. No \"Balance\" field")
	}
	if id, ok := mapBody["Id"]; ok {
		if id != float64(user.Id) {
			clearTable(user.Id, 0)
			t.Errorf("Incorrect user id returned. Expected %d, got %f", user.Id, id)
		}
	} else {
		clearTable(user.Id, 0)
		t.Errorf("Incorrect json output in add_money. No \"Id\" field")
	}
	if statusMsg, ok := mapBody["Status"]; ok {
		if statusMsg != "OK" {
			clearTable(user.Id, 0)
			t.Errorf("Incorrect status returned. Expected %s, got %s", "OK", statusMsg)
		}
	} else {
		clearTable(user.Id, 0)
		t.Errorf("Incorrect json output in add_money. No \"Status\" field")
	}
	clearTable(user.Id, 0)
}

func TestReserveMoney(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.AddMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()

	var user schemes.User
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
	handler = http.HandlerFunc(utils.ReserveMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body = rr.Body.String()
	var mapBody map[string]interface{}
	json.Unmarshal([]byte(body), &mapBody)
	if statusMsg, ok := mapBody["Status"]; ok {
		if statusMsg != "OK" {
			clearTable(user.Id, 1)
			t.Errorf("Incorrect status returned. Expected %s, got %s", "OK", statusMsg)
		}
	} else {
		clearTable(user.Id, 1)
		t.Errorf("Incorrect json output in add_money. No \"Status\" field")
	}
	clearTable(user.Id, 1)
}

func TestTakeMoney(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.AddMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)

	var user schemes.User
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
	handler = http.HandlerFunc(utils.ReserveMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)

	jsonBody = []byte(fmt.Sprintf(`{
		"UserId": %d,
		"OrderId": 1,
		"ServiceId": 1,
		"Cost": 10
	}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/take_money", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(utils.TakeMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	var mapBody map[string]interface{}
	json.Unmarshal([]byte(body), &mapBody)
	if statusMsg, ok := mapBody["Status"]; ok {
		if statusMsg != "OK" {
			clearTable(user.Id, 1)
			t.Errorf("Incorrect status returned. Expected %s, got %s", "OK", statusMsg)
		}
	} else {
		clearTable(user.Id, 1)
		t.Errorf("Incorrect json output in add_money. No \"Status\" field")
	}
	clearTable(user.Id, 1)
}

func TestGetBalance(t *testing.T) {
	jsonBody := []byte(`{"amount": 100}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/add_money", bodyReader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.AddMoney)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)

	var user schemes.User
	json.NewDecoder(rr.Body).Decode(&user)
	jsonBody = []byte(fmt.Sprintf(`{"id": %d}`, user.Id))
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/get_balance", bodyReader)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(utils.GetBalance)
	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	var mapBody map[string]interface{}
	json.Unmarshal([]byte(body), &mapBody)
	if money, ok := mapBody["Balance"]; ok {
		if money != float64(100) {
			clearTable(user.Id, 0)
			t.Errorf("Incorrect balance added. Expected %f, got %f", float32(100), money)
		}
	} else {
		clearTable(user.Id, 0)
		t.Errorf("Incorrect json output in add_money. No \"Balance\" field")
	}
	if id, ok := mapBody["Id"]; ok {
		if id != float64(user.Id) {
			clearTable(user.Id, 0)
			t.Errorf("Incorrect user id returned. Expected %d, got %f", user.Id, id)
		}
	} else {
		clearTable(user.Id, 0)
		t.Errorf("Incorrect json output in add_money. No \"Id\" field")
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
		stmt, err := db.Db.Prepare("DELETE FROM orders WHERE order_id=?")
		if err != nil {
			fmt.Println(err)
		}
		_, err = stmt.Exec(order_id)
		if err != nil {
			fmt.Printf("Could not delete from orders with order_id=%d\n", order_id)
		}
	}
	if user_id > 0 {
		stmt, err := db.Db.Prepare("DELETE FROM user_balances WHERE id=?")
		if err != nil {
			fmt.Println(err)
		}
		_, err = stmt.Exec(order_id)
		if err != nil {
			fmt.Printf("Could not delete from user_balances with user_id=%d\n", user_id)
		}
	}
}
