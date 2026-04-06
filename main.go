package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	apiKey := os.Getenv("HEVY_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: HEVY_API_KEY environment variable is not set.")
		os.Exit(1)
	}

	db, err := InitDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	startContinuousLatestSync(ctx, db)
}

func startContinuousLatestSync(ctx context.Context, db *sql.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	if err := syncLatest(db); err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := syncLatest(db); err != nil {
				fmt.Println(err)
			}
		}
	}
}
