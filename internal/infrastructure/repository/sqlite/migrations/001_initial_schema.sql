-- 筋トレセッションテーブル
CREATE TABLE IF NOT EXISTS strength_trainings (
    id TEXT PRIMARY KEY,
    date DATETIME NOT NULL,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- エクササイズテーブル
CREATE TABLE IF NOT EXISTS exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    training_id TEXT NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    exercise_order INTEGER NOT NULL,
    FOREIGN KEY (training_id) REFERENCES strength_trainings(id) ON DELETE CASCADE
);

-- セットテーブル
CREATE TABLE IF NOT EXISTS sets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise_id INTEGER NOT NULL,
    weight_kg REAL NOT NULL,
    reps INTEGER NOT NULL,
    rest_time_seconds INTEGER NOT NULL,
    rpe INTEGER NULL,
    set_order INTEGER NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_strength_trainings_date ON strength_trainings(date);
CREATE INDEX IF NOT EXISTS idx_exercises_training_id ON exercises(training_id);
CREATE INDEX IF NOT EXISTS idx_exercises_name ON exercises(name);
CREATE INDEX IF NOT EXISTS idx_sets_exercise_id ON sets(exercise_id);

-- 分析用ビュー（パフォーマンス向上）
CREATE VIEW IF NOT EXISTS exercise_max_weights AS
SELECT 
    e.training_id,
    e.name as exercise_name,
    MAX(s.weight_kg) as max_weight,
    st.date
FROM exercises e
JOIN sets s ON e.id = s.exercise_id
JOIN strength_trainings st ON e.training_id = st.id
GROUP BY e.training_id, e.name;

CREATE VIEW IF NOT EXISTS exercise_volumes AS
SELECT 
    e.training_id,
    e.name as exercise_name,
    SUM(s.weight_kg * s.reps) as total_volume,
    st.date
FROM exercises e
JOIN sets s ON e.id = s.exercise_id
JOIN strength_trainings st ON e.training_id = st.id
GROUP BY e.training_id, e.name;
