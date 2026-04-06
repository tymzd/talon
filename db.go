package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

// Bake the database schema CREATE TABLE directives into the compiled binary.
//
//go:embed schema.sql
var schema string

func InitDB() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./talon_test.db"
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	return db, nil
}

// UpsertWorkouts upserts a list of workouts into the database. If a workout
// already exists in the database, its exercises and sets are deleted and the
// new workout's exercises and sets will be inserted.
func UpsertWorkouts(db *sql.DB, workouts []Workout) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Undo the transaction if something goes wrong.
	defer tx.Rollback()

	// Prepare all the insertion SQL statements.
	insertWorkoutStmt, err := tx.Prepare(`
		INSERT INTO workouts (id, title, routine_id, description, start_time, end_time, updated_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET 
			title = excluded.title,
			description = excluded.description,
			end_time = excluded.end_time,
			updated_at = excluded.updated_at;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare workout insert statement: %w", err)
	}
	defer insertWorkoutStmt.Close()

	insertExerciseStmt, err := tx.Prepare(`INSERT INTO exercises (workout_id, sort_order, title, notes) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare exercise insert statement: %w", err)
	}
	defer insertExerciseStmt.Close()

	insertSetStmt, err := tx.Prepare(`INSERT INTO sets (exercise_id, sort_order, set_type, weight_kg, reps, distance_meters, duration_seconds) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare set insert statement: %w", err)
	}
	defer insertSetStmt.Close()

	clearExistingExercisesStmt, err := tx.Prepare(`DELETE FROM exercises WHERE workout_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare clear existing exercises statement: %w", err)
	}
	defer clearExistingExercisesStmt.Close()

	// Insert all workouts, exercises, and sets in the same transaction.
	for _, w := range workouts {
		// Upsert the workout.
		_, err = insertWorkoutStmt.Exec(w.ID, w.Title, w.RoutineID, w.Description, w.StartTime, w.EndTime, w.UpdatedAt, w.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to upsert workout: %w", err)
		}

		// Upsert the exercises and sets. Note that we must delete exercises for
		// this workout if they exist already to avoid duplicates/conflicts if this
		// workout was previously uploaded before.
		_, err = clearExistingExercisesStmt.Exec(w.ID)
		if err != nil {
			return fmt.Errorf("failed to clear old exercises: %w", err)
		}

		for i, ex := range w.Exercises {
			res, err := insertExerciseStmt.Exec(w.ID, i, ex.Title, ex.Notes)
			if err != nil {
				return fmt.Errorf("failed to insert exercise: %w", err)
			}
			exerciseID, _ := res.LastInsertId()
			for j, s := range ex.Sets {
				_, err = insertSetStmt.Exec(exerciseID, j, s.Type, s.WeightKG, s.Reps, s.DistanceMeters, s.DurationSeconds)
				if err != nil {
					return fmt.Errorf("failed to insert set: %w", err)
				}
			}
		}
	}

	return tx.Commit()
}

func MarkLastSynced(db *sql.DB, lastSynced time.Time) error {
	_, err := db.Exec(`INSERT INTO sync_status (id, last_synced_at) VALUES (1, ?) ON CONFLICT (id) DO UPDATE SET last_synced_at = excluded.last_synced_at`, lastSynced)
	if err != nil {
		return fmt.Errorf("failed to mark last synced: %w", err)
	}
	return nil
}
