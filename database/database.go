package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func StartDB() (*sql.DB, error) {
	fmt.Printf("StartDB firing\n")
	connStr := "user=michaelboegner dbname=interviewerio sslmode=disable"

	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database successfully!")

	return db, nil

}
