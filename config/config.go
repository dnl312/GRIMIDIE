package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func ConnectDB() *sql.DB {
	var err error

	// Supabase connection string
	connStr := "postgres://postgres.svkqrngeyxdqkwwyswkc:tl2Vn9nnThKTOpCV@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres"

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print("Error connecting to the database: ", err)
		log.Fatal(err)
	}

	// Check the connection
	if err = DB.Ping(); err != nil {
		log.Print("Error pinging the database: ", err)
		log.Fatal(err)
	}

	log.Print("Connected to the Supabase database")

	return DB
}
