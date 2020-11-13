package main

import (
	"microserviceDemo/database"
	"microserviceDemo/person"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	database.SetupDatabase()
	http.HandleFunc("/people", person.PersonHandler)
	http.HandleFunc("/people/", person.SinglePersonHandler)
	http.ListenAndServe(":5000", nil)

}
