package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB connects to the database
func ConnectDB() *sql.DB {
	var err error
	DB, err := sql.Open("mysql", "root:@tcp(localhost:3306)/GRIMIDIE")
	if err != nil {
		log.Print("Error connecting to the database: ", err)
		log.Fatal(err)
	}

	// Check the connection
	if err = DB.Ping(); err != nil {
		log.Print("Error pinging the database: ", err)
		log.Fatal(err)
	}

	log.Print("Connected to the database")

	return DB
}
