package main

import "time"

type PaginatedWorkoutsResponse struct {
	Page      int     `json:"page"`
	PageCount int     `json:"page_count"`
	Events    []Event `json:"events"`
}

type Event struct {
	Type    string  `json:"type"`
	Workout Workout `json:"workout"`
}

type Workout struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	RoutineID   string     `json:"routine_id"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Exercises   []Exercise `json:"exercises"`
}

type Exercise struct {
	Index int    `json:"index"`
	Title string `json:"title"`
	Notes string `json:"notes"`
	Sets  []Set  `json:"sets"`
}

type Set struct {
	Index           int      `json:"index"`
	Type            SetType  `json:"type"`
	WeightKG        *float64 `json:"weight_kg"`
	Reps            *int     `json:"reps"`
	DistanceMeters  *float64 `json:"distance_meters"`
	DurationSeconds *int     `json:"duration_seconds"`
	RPE             *float64 `json:"rpe"`
	CustomMetric    *float64 `json:"custom_metric"`
}

type SetType string

const (
	SetTypeWarmup  SetType = "warmup"
	SetTypeNormal  SetType = "normal"
	SetTypeFailure SetType = "failure"
	SetTypeDrop    SetType = "dropset"
)
