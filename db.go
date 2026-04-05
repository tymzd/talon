package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./talon_test.db")
	if err != nil {
		return nil, err
	}
	fmt.Println("Initialised!")
	return db, nil
}
