package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

var hevyBaseURL = "https://api.hevyapp.com/v1"
var limiter = rate.NewLimiter(rate.Every(1*time.Second), 1)

func getWorkoutsSince(ctx context.Context, logger *slog.Logger, apiKey string, since time.Time) ([]Workout, error) {
	params := url.Values{}
	params.Add("pageSize", "10")
	params.Add("since", since.Format(time.RFC3339))

	var workouts []Workout
	pageCount := 1
	for i := 1; i <= pageCount; i++ {
		params.Set("page", strconv.Itoa(i))
		fullURL := fmt.Sprintf("%s/workouts/events?%s", hevyBaseURL, params.Encode())
		logger.Info(fmt.Sprintf("Fetching page %d: %q", i, fullURL))

		req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		req.Header.Set("accept", "application/json")
		req.Header.Set("api-key", apiKey)

		// Wait for the rate limiter to allow the request.
		if err := limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("error waiting for rate limiter: %w", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %w", err)
		}
		if err := func() error {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response: %w", err)
			}

			var respData PaginatedWorkoutsResponse
			if err := json.Unmarshal(body, &respData); err != nil {
				return fmt.Errorf("failed to unmarshal workouts: %w", err)
			}
			pageCount = respData.PageCount
			for _, w := range respData.Events {
				workouts = append(workouts, w.Workout)
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}
	return workouts, nil
}
