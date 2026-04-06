package main

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// setupTestDB creates a new in-memory SQLite database and initialises the
// schema.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory sqlite db: %v", err)
	}
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}
	return db
}

func TestUpsertWorkouts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert 1 workout.
	exampleWorkout := Workout{
		ID:        "test-uuid-1",
		Title:     "Leg Day",
		RoutineID: "routine-uuid-1",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Exercises: []Exercise{
			{
				Title: "Squat",
				Sets: []Set{
					{Type: SetTypeNormal, WeightKG: ptr(100.0), Reps: ptr(5)},
				},
			},
		},
	}
	if err := UpsertWorkouts(db, []Workout{exampleWorkout}); err != nil {
		t.Fatalf("UpsertWorkouts failed: %v", err)
	}

	// Assert that 1 workout and 1 set was inserted.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM workouts WHERE id = ?", exampleWorkout.ID).Scan(&count); err != nil || count != 1 {
		t.Errorf("Expected 1 workout in DB, got %d", count)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM sets").Scan(&count); err != nil || count != 1 {
		t.Errorf("Expected 1 set in DB, got %d", count)
	}

	// Modify the existing workout,
	exampleWorkout.Title = "Crazy Leg Day"
	exampleWorkout.EndTime = exampleWorkout.EndTime.Add(time.Hour)
	exampleWorkout.Exercises = []Exercise{
		{
			Title: "Squat",
			Sets: []Set{
				{Type: SetTypeNormal, WeightKG: ptr(100.0), Reps: ptr(5)},
				{Type: SetTypeNormal, WeightKG: ptr(100.0), Reps: ptr(5)},
				{Type: SetTypeNormal, WeightKG: ptr(100.0), Reps: ptr(5)},
				{Type: SetTypeNormal, WeightKG: ptr(100.0), Reps: ptr(5)},
				{Type: SetTypeNormal, WeightKG: ptr(100.0), Reps: ptr(5)},
			},
		},
		{
			Title: "Bulgarian Split Squat",
			Sets: []Set{
				{Type: SetTypeNormal, WeightKG: ptr(10.0), Reps: ptr(20)},
				{Type: SetTypeNormal, WeightKG: ptr(10.0), Reps: ptr(15)},
				{Type: SetTypeNormal, WeightKG: ptr(10.0), Reps: ptr(10)},
			},
		},
	}
	if err := UpsertWorkouts(db, []Workout{exampleWorkout}); err != nil {
		t.Fatalf("UpsertWorkouts failed: %v", err)
	}

	// Assert that 1 workout and 8 sets were inserted.
	if err := db.QueryRow("SELECT COUNT(*) FROM workouts WHERE id = ?", exampleWorkout.ID).Scan(&count); err != nil || count != 1 {
		t.Errorf("Expected 1 workout in DB, got %d", count)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM sets").Scan(&count); err != nil || count != 8 {
		t.Errorf("Expected 8 sets in DB, got %d", count)
	}
}

func ptr[T any](v T) *T {
	return &v
}
