package main

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"
)

// Bake the database schema CREATE TABLE directives into the compiled binary.
//
//go:embed schema.sql
var schema string

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./talon_test.db")
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	return db, nil
}
