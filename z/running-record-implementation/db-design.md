# データベース設計詳細

## 1. テーブル設計

### running_sessions テーブル

```sql
CREATE TABLE IF NOT EXISTS running_sessions (
    id TEXT PRIMARY KEY,              -- セッションID（UUID）
    date DATETIME NOT NULL,           -- ランニング日
    distance_km REAL NOT NULL,        -- 距離（km）
    duration_seconds INTEGER NOT NULL, -- 時間（秒）
    pace_seconds_per_km REAL NOT NULL, -- ペース（秒/km）
    heart_rate_bpm INTEGER NULL,      -- 心拍数（オプション）
    run_type TEXT NOT NULL,           -- ランニングタイプ
    notes TEXT,                       -- メモ
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### カラム詳細

| カラム名 | 型 | NULL | 説明 | 例 |
|---------|---|------|------|---|
| id | TEXT | NOT NULL | UUID形式のセッションID | `550e8400-e29b-41d4-a716-446655440000` |
| date | DATETIME | NOT NULL | ランニング実施日 | `2025-06-16 06:30:00` |
| distance_km | REAL | NOT NULL | 走行距離（km） | `5.0`, `10.5` |
| duration_seconds | INTEGER | NOT NULL | 走行時間（秒） | `1530` (25分30秒) |
| pace_seconds_per_km | REAL | NOT NULL | ペース（秒/km） | `306.0` (5:06/km) |
| heart_rate_bpm | INTEGER | NULL | 平均心拍数 | `165`, `NULL` |
| run_type | TEXT | NOT NULL | ランニングタイプ | `Easy`, `Tempo`, `Interval`, `Long`, `Race` |
| notes | TEXT | NULL | メモ・備考 | `気持ちよく走れた` |
| created_at | DATETIME | NOT NULL | 作成日時 | `2025-06-16 07:00:00` |
| updated_at | DATETIME | NOT NULL | 更新日時 | `2025-06-16 07:00:00` |

## 2. インデックス設計

### パフォーマンス最適化

```sql
-- 日付での検索最適化（期間指定、月次集計等）
CREATE INDEX IF NOT EXISTS idx_running_sessions_date ON running_sessions(date);

-- ランニングタイプでの分析最適化
CREATE INDEX IF NOT EXISTS idx_running_sessions_run_type ON running_sessions(run_type);

-- 距離での検索最適化（距離別分析等）
CREATE INDEX IF NOT EXISTS idx_running_sessions_distance ON running_sessions(distance_km);

-- 複合インデックス（日付 + タイプでの高速検索）
CREATE INDEX IF NOT EXISTS idx_running_sessions_date_type ON running_sessions(date, run_type);
```

### インデックス使用例

```sql
-- 期間指定での検索（idx_running_sessions_date使用）
SELECT * FROM running_sessions 
WHERE date BETWEEN '2025-06-01' AND '2025-06-30';

-- タイプ別分析（idx_running_sessions_run_type使用）
SELECT run_type, COUNT(*), AVG(distance_km) 
FROM running_sessions 
GROUP BY run_type;

-- 長距離ランの検索（idx_running_sessions_distance使用）
SELECT * FROM running_sessions 
WHERE distance_km >= 10.0;
```

## 3. 分析用ビュー

### 週次統計ビュー

```sql
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
```

### 月次統計ビュー

```sql
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
```

### ランニングタイプ別統計ビュー

```sql
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
```

## 4. 制約とバリデーション

### データベースレベルの制約

```sql
-- 距離は正の値のみ
CHECK (distance_km > 0)

-- 時間は正の値のみ
CHECK (duration_seconds > 0)

-- ペースは正の値のみ
CHECK (pace_seconds_per_km > 0)

-- 心拍数は正の値（NULLは許可）
CHECK (heart_rate_bpm IS NULL OR heart_rate_bpm > 0)

-- ランニングタイプは定義済み値のみ
CHECK (run_type IN ('Easy', 'Tempo', 'Interval', 'Long', 'Race'))
```

### アプリケーションレベルのバリデーション

```go
// 距離のバリデーション
func ValidateDistance(km float64) error {
    if km <= 0 {
        return errors.New("距離は正の値である必要があります")
    }
    return nil
}

// 時間のバリデーション
func ValidateDuration(seconds int) error {
    if seconds <= 0 {
        return errors.New("時間は正の値である必要があります")
    }
    return nil
}

// ランニングタイプのバリデーション
func ValidateRunType(runType string) error {
    validTypes := []string{"Easy", "Tempo", "Interval", "Long", "Race"}
    for _, valid := range validTypes {
        if runType == valid {
            return nil
        }
    }
    return fmt.Errorf("無効なランニングタイプ: %s", runType)
}
```

## 5. 既存テーブルとの関係

### 独立性の確保

```sql
-- 既存の筋トレテーブル
strength_trainings (id, date, notes, created_at, updated_at)
exercises (id, training_id, name, exercise_order)
sets (id, exercise_id, weight_kg, reps, rpe, set_order)

-- 新しいランニングテーブル
running_sessions (id, date, distance_km, duration_seconds, ...)

-- 外部キー関係なし - 完全に独立
```

### 共通要素

- **ID形式**: 両方ともUUIDを使用
- **日付形式**: 両方ともDATETIME型
- **メモ機能**: 両方ともnotesカラム
- **タイムスタンプ**: 両方ともcreated_at, updated_at

## 6. マイグレーション戦略

### 段階的マイグレーション

```sql
-- 003_add_running_tables.sql
-- Phase 1: テーブル作成
CREATE TABLE IF NOT EXISTS running_sessions (...);

-- Phase 2: インデックス作成
CREATE INDEX IF NOT EXISTS ...;

-- Phase 3: ビュー作成
CREATE VIEW IF NOT EXISTS ...;

-- Phase 4: 制約追加（必要に応じて）
-- ALTER TABLE running_sessions ADD CONSTRAINT ...;
```

### ロールバック計画

```sql
-- ロールバック用SQL（必要時）
DROP VIEW IF EXISTS running_weekly_stats;
DROP VIEW IF EXISTS running_monthly_stats;
DROP VIEW IF EXISTS running_type_stats;
DROP INDEX IF EXISTS idx_running_sessions_date;
DROP INDEX IF EXISTS idx_running_sessions_run_type;
DROP INDEX IF EXISTS idx_running_sessions_distance;
DROP INDEX IF EXISTS idx_running_sessions_date_type;
DROP TABLE IF EXISTS running_sessions;
```

## 7. データサンプル

### テストデータ例

```sql
INSERT INTO running_sessions VALUES 
('550e8400-e29b-41d4-a716-446655440001', '2025-06-16 06:30:00', 5.0, 1530, 306.0, 165, 'Easy', '気持ちよく走れた', datetime('now'), datetime('now')),
('550e8400-e29b-41d4-a716-446655440002', '2025-06-17 06:00:00', 8.0, 2112, 264.0, 175, 'Tempo', '息が上がったが維持できた', datetime('now'), datetime('now')),
('550e8400-e29b-41d4-a716-446655440003', '2025-06-18 07:00:00', 3.0, 900, 300.0, NULL, 'Easy', 'リカバリーラン', datetime('now'), datetime('now'));
```

---

この設計により、筋トレ機能と完全に独立したランニング記録システムを構築できます。
