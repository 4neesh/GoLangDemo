package main

import (
	"advanced-web-service/dataStore"
	"advanced-web-service/person"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	dataStore.SetUpDatabase()

	personHandler := http.HandlerFunc(person.PersonHandler)
	singlePersonHandler := http.HandlerFunc(person.SinglePersonHandler)
	http.Handle("/people", middlewareHandler(personHandler))
	http.Handle("/people/", middlewareHandler(singlePersonHandler))
	http.ListenAndServe(":5000", nil)

}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware")
		handler.ServeHTTP(w, r)

	})
}
