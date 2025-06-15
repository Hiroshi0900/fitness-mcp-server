# T1-1: マイグレーションファイル作成 - テストケース洗い出し

## 対象
- ファイル: `internal/infrastructure/repository/sqlite/migrations/003_add_running_tables.sql`
- 目的: ランニング専用テーブル・インデックス・ビューの作成

## テストケース一覧

### 1. スキーマ作成テスト

#### TC1-1: running_sessions テーブル作成
```go
func TestMigration003_CreateRunningSessionsTable(t *testing.T) {
    // Arrange: テスト用DB
    // Act: マイグレーション実行  
    // Assert: テーブルが存在し、正しいカラム構成を持つ
}
```

**検証項目**:
- テーブル存在確認
- カラム名・型・制約の確認
- PRIMARY KEY設定確認
- NOT NULL制約確認
- DEFAULT値確認

#### TC1-2: カラム型検証
```go
func TestMigration003_ColumnTypes(t *testing.T) {
    // Assert: 各カラムの型が正しい
    // id: TEXT
    // date: DATETIME  
    // distance_km: REAL
    // duration_seconds: INTEGER
    // pace_seconds_per_km: REAL
    // heart_rate_bpm: INTEGER (NULL許可)
    // run_type: TEXT
    // notes: TEXT (NULL許可)
    // created_at: DATETIME (DEFAULT)
    // updated_at: DATETIME (DEFAULT)
}
```

### 2. インデックス作成テスト

#### TC2-1: 基本インデックス存在確認
```go
func TestMigration003_CreateIndexes(t *testing.T) {
    // Assert: 必要なインデックスが全て作成されている
    // - idx_running_sessions_date
    // - idx_running_sessions_run_type  
    // - idx_running_sessions_distance
}
```

#### TC2-2: インデックス効果確認
```go
func TestMigration003_IndexPerformance(t *testing.T) {
    // Arrange: テストデータ挿入
    // Act: インデックス対象カラムでのクエリ実行
    // Assert: EXPLAINでインデックス使用確認
}
```

### 3. ビュー作成テスト

#### TC3-1: 週次統計ビュー
```go
func TestMigration003_CreateWeeklyStatsView(t *testing.T) {
    // Arrange: 複数週のテストデータ
    // Act: running_weekly_stats ビューにクエリ
    // Assert: 正しい週次集計結果
}
```

#### TC3-2: 月次統計ビュー  
```go
func TestMigration003_CreateMonthlyStatsView(t *testing.T) {
    // Arrange: 複数月のテストデータ
    // Act: running_monthly_stats ビューにクエリ
    // Assert: 正しい月次集計結果
}
```

#### TC3-3: ランニングタイプ別統計ビュー
```go
func TestMigration003_CreateTypeStatsView(t *testing.T) {
    // Arrange: 異なるrun_typeのテストデータ
    // Act: running_type_stats ビューにクエリ  
    // Assert: タイプ別の正しい集計結果
}
```

### 4. 制約テスト

#### TC4-1: CHECK制約動作確認
```go
func TestMigration003_CheckConstraints(t *testing.T) {
    // Act & Assert: 無効データ挿入時のエラー確認
    // - distance_km <= 0
    // - duration_seconds <= 0
    // - pace_seconds_per_km <= 0
    // - heart_rate_bpm <= 0
    // - run_type NOT IN ('Easy', 'Tempo', 'Interval', 'Long', 'Race')
}
```

### 5. 既存機能への影響確認テスト

#### TC5-1: 筋トレテーブル非影響確認
```go
func TestMigration003_ExistingTablesUnaffected(t *testing.T) {
    // Arrange: 既存筋トレデータ
    // Act: マイグレーション実行
    // Assert: 既存データ・構造に変更なし
}
```

#### TC5-2: 筋トレ機能動作確認
```go
func TestMigration003_StrengthTrainingStillWorks(t *testing.T) {
    // Arrange: マイグレーション後のDB
    // Act: 筋トレ記録・取得操作
    // Assert: 正常動作確認
}
```

### 6. ロールバックテスト

#### TC6-1: マイグレーション巻き戻し
```go
func TestMigration003_Rollback(t *testing.T) {
    // Arrange: マイグレーション適用済みDB
    // Act: ロールバック実行
    // Assert: ランニング関連オブジェクトが全て削除
}
```

### 7. データ操作テスト

#### TC7-1: 基本INSERT/SELECT
```go
func TestMigration003_BasicDataOperations(t *testing.T) {
    // Act: サンプルデータ挿入
    // Assert: 正常に挿入・取得できる
}
```

#### TC7-2: 日付フィルタリング
```go
func TestMigration003_DateFiltering(t *testing.T) {
    // Arrange: 異なる日付のデータ
    // Act: 期間指定クエリ
    // Assert: 正しくフィルタリングされる
}
```

## テスト実装計画

### ファイル構成
```
internal/infrastructure/repository/sqlite/
├── migrations/
│   ├── 003_add_running_tables.sql     # 実装対象
│   └── 003_add_running_tables_test.go # テストファイル
```

### テストヘルパー
```go
// テスト用DB作成
func setupTestDBForMigration(t *testing.T) *sql.DB

// マイグレーション実行
func runMigration(db *sql.DB, filename string) error

// テーブル存在確認
func assertTableExists(t *testing.T, db *sql.DB, tableName string)

// インデックス存在確認  
func assertIndexExists(t *testing.T, db *sql.DB, indexName string)

// カラム情報取得
func getColumnInfo(db *sql.DB, tableName string) ([]ColumnInfo, error)
```

### 実行順序
1. **テスト実装** (上記テストケースを全てコメントアウトして作成)
2. **マイグレーションSQL作成** (テストが通るように実装)
3. **リファクタリング** (パフォーマンス・可読性向上)
4. **lint + test 実行** (品質確認)

## 次のアクション

このテストケース洗い出しに問題がなければ、テスト実装 → マイグレーションSQL作成を開始します。

追加すべきテストケースや修正点があればお知らせください。
