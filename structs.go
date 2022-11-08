package main

// User is Interface for user details.
type User struct {
	Id      int
	Balance float32
}

type UpgradeBalance struct {
	Id     int
	Amount float32 `json:"amount"`
}

type OrderReserve struct {
	UserId    int
	OrderId   int
	ServiceId int
	Cost      float32
}
