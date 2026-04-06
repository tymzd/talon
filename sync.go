package main

import (
	"database/sql"
	"fmt"
)

func syncLatest(db *sql.DB) error {
	fmt.Println("Syncing latest...")
	return nil
}

// TODO: Rate-limited full refresh loop.
