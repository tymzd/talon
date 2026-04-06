# OpenClaw Skill: Hevy Workout Tracker

**Description**: This skill provides access to a local SQLite database (`talon.db`) containing the user's synced workout data from the Hevy app. Use this skill to query workout history and act as the user's personal trainer, keeping them accountable and analysing their progression over time.

---

## 1. Database Schema

```sql
CREATE TABLE workouts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    routine_id TEXT NOT NULL,
    description TEXT NOT NULL,
    start_time TEXT NOT NULL,          -- ISO8601 format (e.g. 2026-04-04T02:42:07+00:00)
    end_time TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_id TEXT NOT NULL,          -- Foreign key
    sort_order INTEGER NOT NULL,
    title TEXT NOT NULL,
    notes TEXT NOT NULL
);

CREATE TABLE sets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise_id INTEGER NOT NULL,      -- Foreign key
    sort_order INTEGER NOT NULL,
    set_type TEXT NOT NULL,            -- e.g., 'normal', 'warmup', 'failure', 'drop'
    weight_kg REAL,                    -- Nullable: Weight in KG (if applicable)
    reps INTEGER,                      -- Nullable: Repetitions (if applicable)
    distance_meters REAL,              -- Nullable: Distance (if cardio)
    duration_seconds INTEGER           -- Nullable: Duration (if cardio/holds)
);
```

## 2. Query Examples

### Accountability

Example queries to help check if the user has been showing up consistently (e.g., today, last 7 days).

*Recent workouts in the last 7 days:*
```bash
sqlite3 -json -readonly ~/talon.db "
SELECT title, start_time, end_time, description
FROM workouts
WHERE date(start_time) >= date('now', '-7 days')
ORDER BY start_time DESC;
"
```

*Did the user workout today?*
```bash
# Using 'localtime' modifier to match user's local timezone
sqlite3 -json -readonly ~/talon.db "
SELECT title, time(start_time) as time_started, time(end_time) as time_finished
FROM workouts
WHERE date(start_time) = date('now', 'localtime');
"
```

### Progression & Plateaus

Example query to identify if the user is getting stronger, increasing volume, or stuck at a plateau.

*Strength Progression on a core lift (e.g., "Squat", "Bench Press") over time:*
```bash
# Evaluates the heaviest weight lifted, max reps on heavy days, and an estimated 1RM.
sqlite3 -json -readonly ~/talon.db "
SELECT 
    date(w.start_time) as workout_date,
    MAX(s.weight_kg) as max_weight,
    MAX(s.reps) as max_reps_on_heavy_set,
    -- Simple Epley formula for tracking 1RM strength progression trends
    MAX(s.weight_kg * (1 + (s.reps / 30.0))) as estimated_1rm
FROM workouts w
JOIN exercises e ON w.id = e.workout_id
JOIN sets s ON e.id = s.exercise_id
WHERE e.title LIKE '%Bench Press%'   -- Adjust filter based on requested exercise
  AND s.set_type IN ('normal', 'failure')
  AND s.weight_kg IS NOT NULL
GROUP BY date(w.start_time)
ORDER BY w.start_time ASC;
"
```
