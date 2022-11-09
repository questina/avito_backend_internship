package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/questina/avito_backend_internship/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title          Balance Management Swagger API
// @version        1.0
// @description    Swagger API for Golang Project Balance Management.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.email kristinagurtov@yandex.ru

// @host     localhost:8000
// @BasePath /
func main() {

	LoadEnv()

	fmt.Println("Server will start at http://localhost:8000/")

	connectDatabase()

	route := mux.NewRouter()

	addApproutes(route)

	route.HandleFunc("/swagger/*", httpSwagger.Handler())

	log.Fatal(http.ListenAndServe(":8000", route))
}
