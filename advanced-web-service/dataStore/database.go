package dataStore

import (
	"database/sql"
	"fmt"
	"log"
)

var DbConn *sql.DB

func SetUpDatabase() {
	var err error
	fmt.Println("setting up database")
	DbConn, err := sql.Open("mysql", "root:Sh==na1992@tcp(127.0.0.1:3306)/people_store")
	if err != nil {
		log.Fatal(err)

	}
	if DbConn == nil {
		log.Fatal("Error in database")
	}

}
