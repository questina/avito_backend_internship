package main

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

type AddMoneyReturn struct {
	Id      int
	Balance float32
	Status  string
}

type StatusMessage struct {
	Status string
}
