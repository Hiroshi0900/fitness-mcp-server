-- Remove RestTime and Category fields migration
-- This migration removes the category field from exercises table and rest_time_seconds field from sets table

-- 1. Drop views first to avoid dependency issues
DROP VIEW IF EXISTS exercise_max_weights;
DROP VIEW IF EXISTS exercise_volumes;

-- 2. Remove category column from exercises table
-- SQLite doesn't support DROP COLUMN directly, so we need to recreate the table
CREATE TABLE exercises_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    training_id TEXT NOT NULL,
    name TEXT NOT NULL,
    exercise_order INTEGER NOT NULL,
    FOREIGN KEY (training_id) REFERENCES strength_trainings(id) ON DELETE CASCADE
);

-- Copy data (excluding category column)
INSERT INTO exercises_new (id, training_id, name, exercise_order)
SELECT id, training_id, name, exercise_order FROM exercises;

-- Drop old table and rename new one
DROP TABLE exercises;
ALTER TABLE exercises_new RENAME TO exercises;

-- 3. Remove rest_time_seconds column from sets table
CREATE TABLE sets_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise_id INTEGER NOT NULL,
    weight_kg REAL NOT NULL,
    reps INTEGER NOT NULL,
    rpe INTEGER NULL,
    set_order INTEGER NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

-- Copy data (excluding rest_time_seconds column)
INSERT INTO sets_new (id, exercise_id, weight_kg, reps, rpe, set_order)
SELECT id, exercise_id, weight_kg, reps, rpe, set_order FROM sets;

-- Drop old table and rename new one
DROP TABLE sets;
ALTER TABLE sets_new RENAME TO sets;

-- 4. Recreate indexes for new table structure
CREATE INDEX IF NOT EXISTS idx_exercises_training_id ON exercises(training_id);
CREATE INDEX IF NOT EXISTS idx_exercises_name ON exercises(name);
CREATE INDEX IF NOT EXISTS idx_sets_exercise_id ON sets(exercise_id);

-- 5. Recreate views with updated schema (no category, no rest_time_seconds)
CREATE VIEW exercise_max_weights AS
SELECT 
    e.training_id,
    e.name as exercise_name,
    MAX(s.weight_kg) as max_weight,
    st.date
FROM exercises e
JOIN sets s ON e.id = s.exercise_id
JOIN strength_trainings st ON e.training_id = st.id
GROUP BY e.training_id, e.name;

CREATE VIEW exercise_volumes AS
SELECT 
    e.training_id,
    e.name as exercise_name,
    SUM(s.weight_kg * s.reps) as total_volume,
    st.date
FROM exercises e
JOIN sets s ON e.id = s.exercise_id
JOIN strength_trainings st ON e.training_id = st.id
GROUP BY e.training_id, e.name;
