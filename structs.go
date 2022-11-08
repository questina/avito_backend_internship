package main

// User is Interface for user details.
type User struct {
	Id      int
	Balance int
}

type UpgradeBalance struct {
	Id     int
	Amount int `json:"amount"`
}

type OrderReserve struct {
	UserId    int
	OrderId   int
	ServiceId int
	Cost      int
}

// ErrorResponse is interface for sending error message with code.
type ErrorResponse struct {
	Code    int
	Message string
}
