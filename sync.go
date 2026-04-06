package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func syncLatest(ctx context.Context, db *sql.DB, logger *slog.Logger, apiKey string) error {
	startTime := time.Now()

	// Get the last updated time. If it doesn't exist, default to the UNIX
	// epoch.
	var lastUpdated time.Time
	err := db.QueryRow(`SELECT last_synced_at FROM sync_status WHERE id = 1`).Scan(&lastUpdated)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to retrieve last updated time: %w", err)
	}
	logger.Info(fmt.Sprintf("Syncing latest workouts since %s", lastUpdated.String()))

	// Retrieve workouts from Hevy since the last updated time.
	workouts, err := getWorkoutsSince(ctx, logger, apiKey, lastUpdated)
	if err != nil {
		return fmt.Errorf("failed to retrieve workouts: %w", err)
	}
	if len(workouts) == 0 {
		logger.Info("No new workouts found")
		return nil
	}

	// Commit all workouts to the database and set the last updated time.
	logger.Info(fmt.Sprintf("Upserting %d workouts", len(workouts)))
	if err := UpsertWorkouts(db, workouts); err != nil {
		return fmt.Errorf("failed to commit workouts to database: %w", err)
	}
	if err := MarkLastSynced(db, startTime); err != nil {
		return fmt.Errorf("failed to mark last synced: %w", err)
	}

	logger.Info(fmt.Sprintf("Sync completed in %s", time.Since(startTime).String()))
	return nil
}

// TODO: Rate-limited full refresh loop.
