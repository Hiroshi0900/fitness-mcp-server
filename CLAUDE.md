# Fitness MCP Server - Claude開発ガイド

## プロジェクト概要

これは筋力トレーニングとランニング記録を管理するGoで書かれた**Fitness MCP（Model Context Protocol）Server**です。
このサーバーはMCPプロトコルを実装し、Claude Desktopやその他のMCPクライアントと統合可能なフィットネス追跡ツールを提供します。

## アーキテクチャ

### クリーンアーキテクチャの実装

コードベースは**クリーンアーキテクチャの原則**に従い、関心の分離を明確にしています：

```
internal/
├── domain/           # コアビジネスロジック（エンティティ、値オブジェクト）
├── application/      # ユースケース（コマンド/クエリハンドラー）
├── infrastructure/   # 外部依存関係（SQLiteなど）
└── interface/        # アダプター（MCPツール、リポジトリ）
```

### 主要なアーキテクチャパターン

#### 1. **コマンドクエリ責任分離（CQRS）**
- **コマンド側**: 書き込み操作を処理（トレーニング記録）
- **クエリ側**: 読み込み操作を処理（トレーニング取得、自己記録）
- 各側面で別々のハンドラー、DTO、ユースケース

#### 2. **ドメイン駆動設計（DDD）**
- ビジネスロジックのカプセル化を持つリッチドメインエンティティ
- プリミティブ強迫を避けるための値オブジェクト
- 複雑なビジネスルールのためのドメインサービス

#### 3. **リポジトリパターン**
- インターフェースを通じた抽象化されたデータアクセス
- 組み込みマイグレーション付きSQLite実装
- 読み書き分離リポジトリ

## ドメインモデル

### コアエンティティ

#### **StrengthTraining**（集約ルート）
```go
type StrengthTraining struct {
    id        shared.TrainingID
    date      time.Time
    exercises []*Exercise
    notes     string
}
```

#### **Exercise**（エンティティ）
```go
type Exercise struct {
    name     ExerciseName     // 値オブジェクト
    sets     []Set           // セットのコレクション
    category ExerciseCategory // 値オブジェクト
}
```

#### **Set**（値オブジェクト）
```go
type Set struct {
    weight   Weight   // 値オブジェクト
    reps     Reps     // 値オブジェクト  
    restTime RestTime // 値オブジェクト
    rpe      *RPE     // オプションのRPE評価
}
```

### 値オブジェクト
- **Weight**: 現実的な重量範囲を検証（0-1000kg）
- **Reps**: 反復回数を検証（1-500）
- **RestTime**: 休憩時間を検証（0-30分）
- **RPE**: 自覚的運動強度（1-10スケール）
- **ExerciseName**: 強く型付けされたエクササイズ名
- **ExerciseCategory**: コンパウンド/アイソレーション/有酸素運動の分類

## MCPツール統合

### 利用可能なツール

#### 1. **record_training**
- 複数のエクササイズを含む完全なトレーニングセッションを記録
- エクササイズカテゴリとセットデータの検証
- RPE（自覚的運動強度）追跡をサポート
- 詳細なエラーメッセージ付きの豊富なパラメータ検証

#### 2. **get_trainings_by_date_range**
- 指定された日付範囲内のトレーニングセッションを取得
- 最大1年間の期間制限
- 要約と統計付きの整形されたレスポンス

#### 3. **get_personal_records**
- エクササイズごとの最大重量、最大レップ数、最大ボリュームを計算
- オプションのエクササイズ名フィルタリング
- 日付付きの詳細な達成追跡

### MCPプロトコル実装
- `github.com/mark3labs/mcp-go`ライブラリを使用
- 包括的な説明付きの構造化ツール定義
- ユーザーフレンドリーなメッセージでのエラーハンドリング
- タイムアウト保護（リクエストあたり30秒）

## データベース設計

### SQLiteスキーマ
```sql
-- コアテーブル
strength_trainings (id, date, notes, timestamps)
exercises (id, training_id, name, category, order)
sets (id, exercise_id, weight_kg, reps, rest_time_seconds, rpe, order)

-- パフォーマンスビュー
exercise_max_weights -- 事前計算された最大重量
exercise_volumes     -- 事前計算されたボリューム
```

### マイグレーション戦略
- `//go:embed`を使用した組み込みSQLファイル
- 起動時の自動マイグレーション
- ファイル命名によるバージョン管理

## 開発ワークフロー

### ビルドコマンド（Makefileから）

#### ローカル開発
```bash
make build          # ローカルでバイナリをビルド
make test           # Goテストを実行
make lint           # 静的解析（vet + fmt）
make check          # lintとtestの組み合わせ
make clean          # ビルド成果物をクリーンアップ
```

#### Docker開発
```bash
make docker-build   # Dockerイメージをビルド
make docker-dev     # シェル付きの開発コンテナを開始
make run            # DockerでMCPサーバーを開始
make stop           # MCPサーバーを停止
make logs           # サーバーログを表示
make status         # コンテナステータスを確認
```

#### テスト
```bash
make docker-test    # フルテストスイート（記録 + クエリ）
make test-record    # 記録機能のテスト
make test-query     # クエリ機能のテスト  
make test-local     # ローカルMCP統合テスト
make test-interactive # インタラクティブテストモード
```

#### Claude Desktop統合
```bash
make setup-claude-local   # ローカルバイナリのセットアップ
make setup-claude-docker  # Dockerのセットアップ
make check-claude-config  # 設定の検証
```

#### CI/CD
```bash
make ci             # フルCIパイプライン（lint + test + MCPテスト）
make dev-setup      # 完全な開発環境のセットアップ
```

## コードパターンと規約・実装について
- 基本的なコーディングはGoの慣習やベストプラクティスに沿って実装する
- TDDを推奨する。そのため実装時はテストケースの洗い出しとテストケースをUnitTestにコメントアウトして書き始めることが望ましい。
- テストのケースの洗い出しを一旦提出すること。問題なければプロダクトコードの実装とテストの実装を行うこと
- ITテストはハンドラの粒度で行うことを想定する。ただし対応が拡大になる可能性があるので、TDDのような形式で実装は不要。
- テストケースの洗い出しと実装計画（ファイルの追加や変更箇所の検討）は同時でも良い。
- 不必要にスクリプトやドキュメントを作成しなくていい。
- 実装計画などをファイルとして残す場合は/zディレクトリの中にサブディレクトリを作って実装すること
- 実装を進めた場合は必ずlintとテストを実行し、問題ないことを確認すること

## 設定

### デフォルト動作
- データベースパス: `./data/fitness.db`（実行ファイルからの相対パス）
- 自動ディレクトリ作成
- 必要に応じて一時ディレクトリへの優雅なフォールバック

## 特別な考慮事項

### パフォーマンス
- パフォーマンス最適化されたビュー付きのSQLite
- 接続プール設定
- 一般的なクエリのためのインデックス最適化
- タイムアウト付きのgoroutineベースのリクエスト処理

### データ整合性
- 複数テーブル操作での トランザクション一貫性
- CASCADE削除付きの外部キー制約
- 複数レイヤーでのドメイン検証
- 可能な場合の冪等操作

### MCPプロトコル準拠
- 適切なツール登録とメタデータ
- パラメータ検証とエラーレポート
- 整形付きの構造化レスポンス
- 長時間実行操作のタイムアウトハンドリング

### 開発ツール
- 整理されたコマンド付きの包括的なMakefile
- 一貫した環境のためのDockerサポート
- MCP機能の統合テストスクリプト
- 修正を行なって改修を進める場合は claude code のMCPサーバを利用することを推奨する。
    - Intelijエディタも利用できると思うが、利用する場合は事前に理由も含めて相談してほしい 

## 将来の拡張

アーキテクチャは以下の簡単な拡張をサポートします：
- **ランニングモジュール**: `internal/domain/running/`で既にスキャフォールド済み
- **追加スポーツ**: モジュラードメイン設計
- **分析**: ビュー付きのクエリ側最適化
- **Webインターフェース**: クリーンな分離により複数のインターフェースが可能
- **エクスポート機能**: リポジトリ抽象化により複数のバックエンドをサポート

## 依存関係

### コア
- `github.com/mark3labs/mcp-go` - MCPプロトコル実装
- `modernc.org/sqlite` - Pure Go SQLiteドライバー
- ほとんどの機能で標準ライブラリ

### 開発
- コンテナ化のためのDocker & Docker Compose
- ビルド自動化のためのMake
- 統合テストのためのシェルスクリプト

このアーキテクチャは、関心の分離の明確化、包括的なテスト、
将来の機能の簡単な拡張性を持つフィットネス追跡の堅実な基盤を提供します。