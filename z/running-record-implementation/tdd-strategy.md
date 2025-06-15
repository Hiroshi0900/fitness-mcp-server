# TDD実装戦略

## TDDアプローチ

### 基本方針
1. **Red-Green-Refactor サイクル** を厳密に守る
2. **既存のテスト構造** を踏襲
3. **ドメインファースト** - ビジネスロジックから実装
4. **段階的なテスト追加** - Phase毎にテスト完成度を高める

## Phase別TDD計画

### Phase 1: 基盤整備（データベース）
- **テスト対象**: マイグレーション結果
- **テスト種類**: 
  - スキーマ検証テスト
  - 制約動作テスト
  - インデックス存在確認テスト

```go
// マイグレーションテスト例
func TestMigration003_RunningTables(t *testing.T) {
    // Arrange: テスト用DB
    db := setupTestDB(t)
    
    // Act: マイグレーション実行
    err := runMigration(db, "003_add_running_tables.sql")
    
    // Assert: テーブル・インデックス存在確認
    assert.NoError(t, err)
    assertTableExists(t, db, "running_sessions")
    assertIndexExists(t, db, "idx_running_sessions_date")
}
```

### Phase 2: Repository層
- **Red**: Repositoryインターフェースの失敗テスト
- **Green**: SQLite実装で最小限の動作
- **Refactor**: エラーハンドリング・パフォーマンス改善

```go
// Repository TDD例
func TestRunningRepository_Save(t *testing.T) {
    // Red: まず失敗するテストを書く
    repo := NewSQLiteRunningRepository(testDB)
    session := createTestRunningSession()
    
    err := repo.Save(session)
    assert.NoError(t, err)
    
    // Green: 最小限の実装で通す
    // Refactor: 品質向上
}
```

### Phase 3: Application層（ドメインモデル中心）
- **ドメインロジック先行**: Value Object から実装
- **コマンドハンドラー**: ユースケース駆動

#### 3.1 Value Object TDD

```go
// Distance TDD例 - Red
func TestDistance_NewDistance_ValidInput(t *testing.T) {
    // Act
    distance, err := NewDistance(5.0)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 5.0, distance.Km())
}

func TestDistance_NewDistance_InvalidInput(t *testing.T) {
    // Act
    distance, err := NewDistance(-1.0)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "distance must be positive")
}
```

#### 3.2 RunningSession TDD

```go
// RunningSession TDD例
func TestRunningSession_NewRunningSession(t *testing.T) {
    // Arrange
    id := shared.NewSessionID()
    date := time.Now()
    distance, _ := NewDistance(5.0)
    duration, _ := NewDuration(time.Minute * 25)
    runType := Easy
    
    // Act
    session, err := NewRunningSession(id, date, distance, duration, runType, "テストラン")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 5.0, session.Distance().Km())
    assert.Equal(t, Easy, session.RunType())
    // ペースが自動計算されることを確認
    expectedPace := 5.0 * 60 // 5分/km = 300秒/km
    assert.Equal(t, expectedPace, session.Pace().SecondsPerKm())
}
```

### Phase 4: MCPツール層
- **パラメータパース**: 入力処理テスト
- **エラーハンドリング**: 各種エラーケーステスト
- **統合テスト**: E2Eでの動作確認

```go
// MCPツール TDD例
func TestRunningTool_ParseDuration_MMSSFormat(t *testing.T) {
    // Arrange
    tool := NewRunningToolHandler(mockHandler)
    
    // Act
    duration, err := tool.parseDuration("25:30")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 25*time.Minute+30*time.Second, duration)
}

func TestRunningTool_ParseDuration_MinutesFormat(t *testing.T) {
    // Act
    duration, err := tool.parseDuration(25.5)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 25*time.Minute+30*time.Second, duration)
}
```

## テスト構造の継承

### 既存パターンの踏襲

```go
// ファイル構成
internal/domain/running/
├── session.go          # 実装
├── session_test.go     # テスト
├── valueobjects.go     # 実装
├── valueobjects_test.go # テスト

// テスト命名規則
func TestRunningSession_MethodName(t *testing.T)
func TestDistance_NewDistance_ValidInput(t *testing.T)
```

### テストヘルパーの活用

```go
// test_helpers.go
func createTestRunningSession() *RunningSession {
    id := shared.NewSessionID()
    distance, _ := NewDistance(5.0)
    duration, _ := NewDuration(25 * time.Minute)
    session, _ := NewRunningSession(id, time.Now(), distance, duration, Easy, "テスト")
    return session
}

func setupTestDB(t *testing.T) *sql.DB {
    // テスト用インメモリDB
}
```

## テスト実行戦略

### 段階的実行
```bash
# Phase別テスト実行
go test ./internal/domain/running/...           # Phase 3
go test ./internal/infrastructure/repository/... # Phase 2  
go test ./internal/interface/mcp-tool/...       # Phase 4

# 全体テスト
make test
```

### カバレッジ確認
```bash
go test -cover ./internal/domain/running/...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 品質ゲート

### 各Phaseの完了基準
- **単体テストカバレッジ**: 90%以上
- **統合テスト**: 主要パス全カバー
- **エラーケーステスト**: 各境界値・異常系
- **回帰テスト**: 既存機能への影響なし

### CI/CD統合
```bash
# CI用テスト実行（Makefileに追加予定）
make ci-test-running  # ランニング機能のみ
make ci-test-all      # 全機能テスト
```

## Red-Green-Refactor 実践

### 実装順序（例：Distance値オブジェクト）

1. **Red**: 失敗するテストを書く
```go
func TestDistance_NewDistance_ValidInput(t *testing.T) {
    distance, err := NewDistance(5.0)
    assert.NoError(t, err)
    assert.Equal(t, 5.0, distance.Km())
}
```

2. **Green**: 最小限の実装で通す
```go
type Distance struct { value float64 }
func NewDistance(km float64) (Distance, error) {
    return Distance{value: km}, nil
}
func (d Distance) Km() float64 { return d.value }
```

3. **Refactor**: バリデーション・品質向上
```go
func NewDistance(km float64) (Distance, error) {
    if km <= 0 {
        return Distance{}, fmt.Errorf("distance must be positive: %f", km)
    }
    return Distance{value: km}, nil
}
```

これにより、既存の品質とテスト体制を維持しながら、新機能を安全に追加できます。
