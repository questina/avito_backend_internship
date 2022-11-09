package main

import (
	"fmt"
	"github.com/gorilla/mux"
)

func addApproutes(route *mux.Router) {

	// route.HandleFunc("/users", getUsers).Methods("GET")

	route.HandleFunc("/add_money", addMoney).Methods("POST")

	route.HandleFunc("/reserve_money", reserveMoney).Methods("POST")

	route.HandleFunc("/take_money", takeMoney).Methods("POST")

	route.HandleFunc("/get_balance", getBalance).Methods("POST")

	route.HandleFunc("/free_money", freeMoney).Methods("POST")

	route.HandleFunc("/generate_report", genReport).Methods("POST")

	route.HandleFunc("/balance_info", balanceInfo).Methods("POST")

	fmt.Println("Routes are Loaded.")
}
