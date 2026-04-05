PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;

CREATE TABLE IF NOT EXISTS sync_status (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    last_synced_at TEXT
);

CREATE TABLE IF NOT EXISTS workouts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    routine_id TEXT NOT NULL,
    description TEXT NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_id TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    title TEXT NOT NULL,
    notes TEXT NOT NULL,
    FOREIGN KEY(workout_id) REFERENCES workouts(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_exercises_workout_id ON exercises(workout_id);

CREATE TABLE IF NOT EXISTS sets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise_id INTEGER NOT NULL,
    sort_order INTEGER NOT NULL,
    set_type TEXT NOT NULL,
    weight_kg REAL,
    reps INTEGER,
    distance_meters REAL,
    duration_seconds INTEGER,
    
    FOREIGN KEY(exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_sets_exercise_id ON sets(exercise_id);
