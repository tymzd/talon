package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetWorkoutsSince(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-key") != "mock-api-key" {
			t.Errorf("Expected api-key header, got %s", r.Header.Get("api-key"))
		}

		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch page {
		case "1":
			w.Write([]byte(`{
				"page": 1,
				"page_count": 2,
				"events": [{"workout": {"id": "mock-uuid-1", "title": "Page 1 Workout"}}]
			}`))
		case "2":
			w.Write([]byte(`{
				"page": 2,
				"page_count": 2,
				"events": [{"workout": {"id": "mock-uuid-2", "title": "Page 2 Workout"}}]
			}`))
		default:
			t.Errorf("Requested unexpected page: %s", page)
		}
	}))
	defer fakeServer.Close()

	// Temporarily swap out the base URL to point to the fake Hevy web server.
	originalURL := hevyBaseURL
	hevyBaseURL = fakeServer.URL
	defer func() { hevyBaseURL = originalURL }()

	workouts, err := getWorkoutsSince(context.Background(), slog.Default(), "mock-api-key", time.Now())
	if err != nil {
		t.Fatalf("getWorkoutsSince failed: %v", err)
	}
	if len(workouts) != 2 {
		t.Fatalf("Expected 2 workouts, got %d", len(workouts))
	}
	if workouts[0].Title != "Page 1 Workout" {
		t.Errorf("Expected 'Page 1 Workout', got '%s'", workouts[0].Title)
	}
	if workouts[1].Title != "Page 2 Workout" {
		t.Errorf("Expected 'Page 2 Workout', got '%s'", workouts[1].Title)
	}
}
