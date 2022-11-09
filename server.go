package main

import (
	_ "github.com/questina/avito_backend_internship/docs"
	"github.com/questina/avito_backend_internship/utils"
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

	route := utils.Start()
	log.Fatal(http.ListenAndServe(":8000", route))
}
