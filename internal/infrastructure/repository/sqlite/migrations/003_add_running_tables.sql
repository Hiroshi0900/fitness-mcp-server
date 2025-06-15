-- ランニングセッションテーブル
CREATE TABLE IF NOT EXISTS running_sessions (
    id TEXT PRIMARY KEY,
    date DATETIME NOT NULL,
    distance_km REAL NOT NULL,
    duration_seconds INTEGER NOT NULL,
    pace_seconds_per_km REAL NOT NULL,
    heart_rate_bpm INTEGER NULL,
    run_type TEXT NOT NULL,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- 制約
    CHECK (distance_km > 0),
    CHECK (duration_seconds > 0),
    CHECK (pace_seconds_per_km > 0),
    CHECK (heart_rate_bpm IS NULL OR heart_rate_bpm > 0),
    CHECK (run_type IN ('Easy', 'Tempo', 'Interval', 'Long', 'Race'))
);

-- インデックス（パフォーマンス最適化）
CREATE INDEX IF NOT EXISTS idx_running_sessions_date ON running_sessions(date);
CREATE INDEX IF NOT EXISTS idx_running_sessions_run_type ON running_sessions(run_type);
CREATE INDEX IF NOT EXISTS idx_running_sessions_distance ON running_sessions(distance_km);
CREATE INDEX IF NOT EXISTS idx_running_sessions_date_type ON running_sessions(date, run_type);

-- 分析用ビュー（週次統計）
CREATE VIEW IF NOT EXISTS running_weekly_stats AS
SELECT 
    DATE(date, 'weekday 0', '-6 days') as week_start,
    COUNT(*) as total_runs,
    SUM(distance_km) as total_distance,
    AVG(distance_km) as avg_distance,
    AVG(pace_seconds_per_km) as avg_pace,
    MIN(pace_seconds_per_km) as best_pace
FROM running_sessions
GROUP BY week_start
ORDER BY week_start DESC;

-- 分析用ビュー（月次統計）
CREATE VIEW IF NOT EXISTS running_monthly_stats AS
SELECT 
    strftime('%Y-%m', date) as month,
    COUNT(*) as total_runs,
    SUM(distance_km) as total_distance,
    AVG(distance_km) as avg_distance,
    AVG(pace_seconds_per_km) as avg_pace,
    MIN(pace_seconds_per_km) as best_pace,
    MAX(distance_km) as longest_run
FROM running_sessions
GROUP BY month
ORDER BY month DESC;

-- 分析用ビュー（ランニングタイプ別統計）
CREATE VIEW IF NOT EXISTS running_type_stats AS
SELECT 
    run_type,
    COUNT(*) as total_runs,
    SUM(distance_km) as total_distance,
    AVG(distance_km) as avg_distance,
    AVG(pace_seconds_per_km) as avg_pace,
    MIN(pace_seconds_per_km) as best_pace,
    MAX(distance_km) as longest_distance
FROM running_sessions
GROUP BY run_type
ORDER BY total_runs DESC;
