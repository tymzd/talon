package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func syncLatest(ctx context.Context, db *sql.DB, logger *slog.Logger, apiKey string) error {
	startTime := time.Now()

	// Get the last updated time. If it doesn't exist, default to the UNIX
	// epoch.
	var lastUpdated time.Time
	var lastUpdatedStr string
	err := db.QueryRow(`SELECT last_synced_at FROM sync_status WHERE id = 1`).Scan(&lastUpdatedStr)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to retrieve last updated time: %w", err)
	}
	if err == nil {
		parsed, parseErr := time.Parse(time.RFC3339Nano, lastUpdatedStr)
		if parseErr != nil {
			parsed, parseErr = time.Parse("2006-01-02 15:04:05.999999999-07:00", lastUpdatedStr)
		}
		if parseErr != nil {
			// Older versions of Talon might have stored the exact output of `time.Time.String()`
			// which includes the monotonic clock reading (e.g. `m=+0.058209224`).
			// Let's strip the ` m=` part and parse the rest.
			cleanStr := lastUpdatedStr
			if mIndex := strings.Index(cleanStr, " m="); mIndex != -1 {
				cleanStr = cleanStr[:mIndex]
			}
			parsed, parseErr = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", cleanStr)
		}
		if parseErr != nil {
			return fmt.Errorf("failed to parse last updated time %q: %w", lastUpdatedStr, parseErr)
		}
		lastUpdated = parsed
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

func syncFull(ctx context.Context, db *sql.DB, logger *slog.Logger, apiKey string) error {
	startTime := time.Now()
	logger.Info("Starting full sync")

	// Retrieve workouts from Hevy since the beginning of time.
	workouts, err := getWorkoutsSince(ctx, logger, apiKey, time.Time{})
	if err != nil {
		return fmt.Errorf("failed to retrieve workouts: %w", err)
	}
	if len(workouts) == 0 {
		logger.Info("No workouts found in full sync")
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

	logger.Info(fmt.Sprintf("Full sync completed in %s", time.Since(startTime).String()))
	return nil
}
