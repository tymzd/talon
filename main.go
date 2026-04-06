package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	_ = godotenv.Load()
	apiKey := os.Getenv("HEVY_API_KEY")
	if apiKey == "" {
		logger.Error("HEVY_API_KEY environment variable is not set.")
		os.Exit(1)
	}

	db, err := InitDB()
	if err != nil {
		logger.Error("Failed to initialise database", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	startContinuousSync(ctx, db, logger, apiKey)
}

func startContinuousSync(ctx context.Context, db *sql.DB, logger *slog.Logger, apiKey string) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	if err := syncLatest(ctx, db, logger, apiKey); err != nil {
		logger.Error("Failed to sync latest workouts", slog.Any("error", err))
	}
	hoursElapsed := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hoursElapsed++
			if hoursElapsed%24 == 0 {
				hoursElapsed = 0
				if err := syncFull(ctx, db, logger, apiKey); err != nil {
					logger.Error("Failed to sync full workouts", slog.Any("error", err))
				}
			} else {
				if err := syncLatest(ctx, db, logger, apiKey); err != nil {
					logger.Error("Failed to sync latest workouts", slog.Any("error", err))
				}
			}
		}
	}
}
