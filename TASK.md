# 筋トレ・ランニング記録管理MCPサーバー 設計書

## プロジェクト概要

### 目的
- **筋トレ（BIG3中心）** と **ランニング** の記録管理
- 目標達成をサポートするパーソナルトレーナー的システム
- MCPサーバーとして会話形式でのデータ分析・アドバイス提供

### 目標設定
- **筋トレ**: ベンチプレス100kg到達（現在95kg、体重72kg）
- **ランニング**: ハーフマラソン出場（7/12）、将来的にフルマラソン
- **完成期限**: 2025年8月まで

### 技術コンセプト
- **言語**: Go
- **アーキテクチャ**: クリーンアーキテクチャ
- **永続化**: SQLite（ローカル、無料）
- **開発手法**: TDD
- **設計**: 軽量DDD + 関数型プログラミング要素

## Phase 1 要件（8月完成目標）

### Core Domain

#### 筋トレドメイン
```go
type StrengthTraining struct {
    ID        TrainingID
    Date      time.Time
    Exercises []Exercise
    Notes     string
}

type Exercise struct {
    Name      ExerciseName      // ベンチプレス、スクワット、デッドリフト
    Sets      []Set
    Category  ExerciseCategory  // Compound, Isolation, Cardio
}

type Set struct {
    Weight     Weight    // 95kg
    Reps       int       // 8回
    RestTime   Duration  // 3分
    RPE        *int      // RPE 8-9（主観的運動強度）
}
```

#### ランニングドメイン
```go
type RunningSession struct {
    ID           SessionID
    Date         time.Time
    Distance     Distance     // 5km, 10km
    Duration     Duration     // 25分30秒
    Pace         Pace         // 5:06/km
    HeartRate    *HeartRate   // 平均150bpm（オプション）
    SessionType  RunType      // Easy, Tempo, Interval, Long
    Notes        string
}

type RunningGoal struct {
    ID          GoalID
    EventType   EventType    // ハーフマラソン、5km、10km
    TargetTime  Duration     // 2時間以内
    EventDate   *time.Time   // 7/12（オプション）
    Status      GoalStatus
}
```

### 機能要件

#### 1. 記録管理（CRUD）
- BIG3の重量・レップ・セット記録
- ランニングの距離・タイム記録
- シンプルな検索・絞り込み

#### 2. 進捗可視化
- ベンチプレス100kg到達への進捗
- ランニングペースの改善傾向
- 月次・週次のトレーニング頻度

#### 3. MCPサーバー機能
- 「今月のベンチプレスのMAX教えて」
- 「ランニングのペース改善してる？」
- 「今日のトレーニングメニュー提案して」

## アーキテクチャ設計

### クリーンアーキテクチャ構成
```
┌─────────────────┐
│   MCP Server    │  ← 会話インターフェース
├─────────────────┤
│  Application    │  ← CQRS (Command/Query Handler)
├─────────────────┤  
│    Domain       │  ← Entity + ValueObject
├─────────────────┤
│ Infrastructure  │  ← SQLite Repository
└─────────────────┘
```

### CQRS設計
```go
// Command側（書き込み）
type StrengthCommandHandler struct {
    repo StrengthRepository
}

func (h *StrengthCommandHandler) RecordTraining(cmd RecordTrainingCommand) error
func (h *StrengthCommandHandler) UpdateTraining(cmd UpdateTrainingCommand) error

type RunningCommandHandler struct {
    repo RunningRepository
}

func (h *RunningCommandHandler) RecordSession(cmd RecordSessionCommand) error

// Query側（読み込み）
type StrengthQueryHandler struct {
    readRepo StrengthReadRepository
}

func (h *StrengthQueryHandler) GetProgressAnalysis(query ProgressQuery) ProgressResult
func (h *StrengthQueryHandler) GetPersonalRecords(query PersonalRecordsQuery) PersonalRecordsResult

type RunningQueryHandler struct {
    readRepo RunningReadRepository
}

func (h *RunningQueryHandler) GetRunningAnalytics(query AnalyticsQuery) AnalyticsResult
func (h *RunningQueryHandler) GetGoalProgress(query GoalProgressQuery) GoalProgressResult
```

### Repository パターン
```go
type WorkoutRepository interface {
    Save(workout *Workout) error
    FindByPeriod(start, end time.Time) ([]*Workout, error)
    FindByExercise(exerciseName string) ([]*Workout, error)
}

type RunningRepository interface {
    Save(session *RunningSession) error
    FindByDateRange(start, end time.Time) ([]*RunningSession, error)
    FindPersonalBests() (map[Distance]PersonalBest, error)
}
```

## 技術選択の根拠

### SQLite選択理由
- **コスト**: 完全無料
- **シンプルさ**: ファイルベース、サーバー不要
- **ポータビリティ**: 単一ファイルでのデータ管理
- **Go対応**: 優秀なドライバー群（modernc.org/sqlite）

### イベントソーシング不採用理由
- 筋トレ記録は「状態遷移」より「事実の記録」
- 必要以上に複雑化し、完成リスクが上がる
- シンプルなCRUD操作で十分

### DDD軽量適用
- **採用**: Entity, Value Object, Repository
- **不採用**: 集約、ドメインサービス（過剰）
- **理由**: 複雑なビジネスルールが少ないため

## ディレクトリ構成

```
fitness-mcp-server/
├── cmd/
│   └── server/
│       └── main.go                 # エントリーポイント
├── internal/
│   ├── domain/
│   │   ├── strength/
│   │   │   ├── entity.go           # StrengthTraining, Exercise
│   │   │   ├── valueobject.go      # Weight, Set, ExerciseName
│   │   │   └── repository.go       # StrengthRepository interface
│   │   ├── running/
│   │   │   ├── entity.go           # RunningSession, RunningGoal
│   │   │   ├── valueobject.go      # Distance, Pace, Duration
│   │   │   └── repository.go       # RunningRepository interface
│   │   └── shared/
│   │       ├── id.go               # TrainingID, SessionID
│   │       └── time.go             # 共通時間型
│   ├── application/
│   │   ├── command/
│   │   │   ├── strength/
│   │   │   │   ├── handler.go      # StrengthCommandHandler
│   │   │   │   └── command.go      # RecordTrainingCommand, etc
│   │   │   └── running/
│   │   │       ├── handler.go      # RunningCommandHandler
│   │   │       └── command.go      # RecordSessionCommand, etc
│   │   └── query/
│   │       ├── strength/
│   │       │   ├── handler.go      # StrengthQueryHandler
│   │       │   ├── query.go        # ProgressQuery, etc
│   │       │   └── result.go       # ProgressResult, etc
│   │       └── running/
│   │           ├── handler.go      # RunningQueryHandler
│   │           ├── query.go        # AnalyticsQuery, etc
│   │           └── result.go       # AnalyticsResult, etc
│   ├── infrastructure/
│   │   ├── repository/
│   │   │   ├── sqlite/
│   │   │   │   ├── strength.go     # SQLite StrengthRepository実装
│   │   │   │   ├── running.go      # SQLite RunningRepository実装
│   │   │   │   └── migration.sql   # DDL
│   │   │   └── memory/             # テスト用InMemoryRepository
│   │   │       ├── strength.go
│   │   │       └── running.go
│   │   └── mcp/
│   │       ├── server.go           # MCPServer本体
│   │       ├── strength_tools.go   # 筋トレ用MCPツール
│   │       ├── running_tools.go    # ランニング用MCPツール
│   │       └── types.go            # MCP通信用型定義
│   └── config/
│       └── config.go               # 設定管理
├── test/
│   ├── integration/                # 統合テスト
│   └── testdata/                   # テスト用データ
├── docs/
│   ├── architecture.md             # アーキテクチャ説明
│   └── api.md                      # MCP API仕様
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```