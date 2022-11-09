package utils

import (
	"fmt"
	"github.com/gorilla/mux"
)

func AddApproutes(route *mux.Router) {

	route.HandleFunc("/add_money", AddMoney).Methods("POST")

	route.HandleFunc("/reserve_money", ReserveMoney).Methods("POST")

	route.HandleFunc("/take_money", TakeMoney).Methods("POST")

	route.HandleFunc("/get_balance", GetBalance).Methods("POST")

	route.HandleFunc("/free_money", FreeMoney).Methods("POST")

	route.HandleFunc("/generate_report", GenReport).Methods("POST")

	route.HandleFunc("/balance_info", BalanceInfo).Methods("POST")

	fmt.Println("Routes are Loaded.")
}
