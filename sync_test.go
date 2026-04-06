package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSyncLatest_InitialSync(t *testing.T) {
	// Set up a fake database and fake Hevy API server.
	db := setupTestDB(t)
	defer db.Close()
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page": 1,
			"page_count": 1,
			"events": [{"workout": {"id": "sync-uuid-1", "title": "Synced Workout"}}]
		}`))
	}))
	defer fakeServer.Close()

	// Temporarily swap out the Hevy API base URL.
	originalURL := hevyBaseURL
	hevyBaseURL = fakeServer.URL
	defer func() { hevyBaseURL = originalURL }()

	logger := slog.Default()
	if err := syncLatest(context.Background(), db, logger, "mock-key"); err != nil {
		t.Fatalf("syncLatest failed: %v", err)
	}

	// Verify that the data on Hevy servers made it into the database.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM workouts WHERE id = 'sync-uuid-1'").Scan(&count); err != nil || count != 1 {
		t.Errorf("Expected the mock workout to be upserted to the DB, but found %d", count)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM sync_status").Scan(&count); err != nil || count != 1 {
		t.Errorf("Expected sync_status to have exactly 1 cursor row, got %d", count)
	}
}

func TestSyncLatest_SubstantialSync(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch page {
		case "1":
			w.Write([]byte(`{
				"page": 1,
				"page_count": 3,
				"events": [
					{
						"type": "updated",
						"workout": {
							"id": "w-uuid-1",
							"title": "Pull Day",
							"routine_id": "r-uuid-1",
							"description": "Standard pull session",
							"start_time": "2026-04-06T09:00:00Z",
							"end_time": "2026-04-06T10:00:00Z",
							"updated_at": "2026-04-06T10:05:00Z",
							"created_at": "2026-04-06T09:00:00Z",
							"exercises": [
								{
									"index": 0,
									"title": "Pull-ups",
									"notes": "Bodyweight",
									"sets": [
										{"index": 0, "type": "normal", "weight_kg": null, "reps": 8, "distance_meters": null, "duration_seconds": null},
										{"index": 1, "type": "normal", "weight_kg": null, "reps": 8, "distance_meters": null, "duration_seconds": null}
									]
								},
								{
									"index": 1,
									"title": "Barbell Row",
									"notes": "Strict form",
									"sets": [
										{"index": 0, "type": "normal", "weight_kg": 60.0, "reps": 10, "distance_meters": null, "duration_seconds": null},
										{"index": 1, "type": "normal", "weight_kg": 60.0, "reps": 10, "distance_meters": null, "duration_seconds": null}
									]
								}
							]
						}
					}
				]
			}`))
		case "2":
			w.Write([]byte(`{
				"page": 2,
				"page_count": 3,
				"events": [
					{
						"type": "updated",
						"workout": {
							"id": "w-uuid-2",
							"title": "Push Day",
							"routine_id": "r-uuid-2",
							"description": "",
							"start_time": "2026-04-05T09:00:00Z",
							"end_time": "2026-04-05T10:00:00Z",
							"updated_at": "2026-04-05T10:05:00Z",
							"created_at": "2026-04-05T09:00:00Z",
							"exercises": [
								{
									"index": 0,
									"title": "Bench Press",
									"notes": "",
									"sets": [
										{"index": 0, "type": "normal", "weight_kg": 80.0, "reps": 5, "distance_meters": null, "duration_seconds": null},
										{"index": 1, "type": "normal", "weight_kg": 80.0, "reps": 5, "distance_meters": null, "duration_seconds": null},
										{"index": 2, "type": "normal", "weight_kg": 80.0, "reps": 5, "distance_meters": null, "duration_seconds": null}
									]
								}
							]
						}
					}
				]
			}`))
		case "3":
			w.Write([]byte(`{
				"page": 3,
				"page_count": 3,
				"events": [
					{
						"type": "updated",
						"workout": {
							"id": "w-uuid-3",
							"title": "Leg Day",
							"routine_id": "r-uuid-3",
							"description": "Heavy squats",
							"start_time": "2026-04-04T09:00:00Z",
							"end_time": "2026-04-04T10:00:00Z",
							"updated_at": "2026-04-04T10:05:00Z",
							"created_at": "2026-04-04T09:00:00Z",
							"exercises": [
								{
									"index": 0,
									"title": "Squat",
									"notes": "",
									"sets": [
										{"index": 0, "type": "normal", "weight_kg": 100.0, "reps": 5, "distance_meters": null, "duration_seconds": null},
										{"index": 1, "type": "failure", "weight_kg": 100.0, "reps": 3, "distance_meters": null, "duration_seconds": null}
									]
								}
							]
						}
					}
				]
			}`))
		default:
			t.Errorf("Requested unexpected page: %s", page)
		}
	}))
	defer fakeServer.Close()

	originalURL := hevyBaseURL
	hevyBaseURL = fakeServer.URL
	defer func() { hevyBaseURL = originalURL }()

	logger := slog.Default()
	if err := syncLatest(context.Background(), db, logger, "mock-key"); err != nil {
		t.Fatalf("syncLatest failed: %v", err)
	}

	// Verify all 3 pages concatenated seamlessly into the DB
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM workouts").Scan(&count); err != nil || count != 3 {
		t.Errorf("Expected 3 workouts in DB, got %d", count)
	}

	if err := db.QueryRow("SELECT COUNT(*) FROM exercises").Scan(&count); err != nil || count != 4 {
		t.Errorf("Expected 4 exercises in DB, got %d", count)
	}

	if err := db.QueryRow("SELECT COUNT(*) FROM sets").Scan(&count); err != nil || count != 9 {
		t.Errorf("Expected 9 sets in DB, got %d", count)
	}

	// Assert we can successfully write complex analytical SQL against our new normalized structure
	var maxWeight float64
	if err := db.QueryRow(`
		SELECT MAX(weight_kg) FROM sets 
		JOIN exercises ON sets.exercise_id = exercises.id 
		WHERE exercises.title = 'Bench Press'
	`).Scan(&maxWeight); err != nil || maxWeight != 80.0 {
		t.Errorf("Failed relational query test! Expected query to pull 80.0kg max for Bench Press, got %v", maxWeight)
	}
}
