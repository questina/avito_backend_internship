package utils

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/questina/avito_backend_internship/db"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Start() *mux.Router {
	LoadEnv()

	fmt.Println("Server will start at http://localhost:8000/")

	db.ConnectDatabase()

	route := mux.NewRouter()

	AddApproutes(route)

	route.HandleFunc("/swagger/*", httpSwagger.Handler())

	return route
}
