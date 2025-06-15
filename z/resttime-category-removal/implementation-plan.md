# RestTime・Category削除対応実装計画

## 概要

トレーニング記録から休憩時間(RestTime)とエクササイズカテゴリ(ExerciseCategory)を削除し、
より直感的で使いやすいインターフェースに変更する。

## 実装状況

### ✅ 完了項目

#### 1. **調査・計画フェーズ**
- [x] 現在のドメインモデルとデータベーススキーマの確認
- [x] RestTimeとExerciseCategoryを削除する影響範囲の調査
- [x] テストケースの洗い出し

#### 2. **ドメインモデル修正**
- [x] `Set`構造体からRestTimeフィールドを削除
  - `RestTime`型と関連メソッドの削除
  - `NewSet`関数の引数変更: `(weight, reps, restTime, rpe)` → `(weight, reps, rpe)`
  - `Set.String()`メソッドから休憩時間表示を削除
- [x] Set関連のテスト修正とテーブルドリブンテスト化
  - Given/When/Then形式への変更
  - RPEありとRPEなしの両パターンのテスト追加

### 🔄 実装中項目

なし

### ❌ 未実装項目

#### 1. **ドメインモデル修正（継続）**
- [ ] `Exercise`構造体からExerciseCategoryフィールドを削除
  - `ExerciseCategory`型と関連メソッドの削除
  - `NewExercise`関数の引数変更: `(name, category)` → `(name)`
  - `Exercise.String()`メソッドからカテゴリ表示を削除
- [ ] Exercise関連のテスト修正とテーブルドリブンテスト化

#### 2. **データベーススキーマ修正**
- [ ] 新しいマイグレーションファイル作成 (`002_remove_resttime_category.sql`)
  - `exercises`テーブルから`category`カラムを削除
  - `sets`テーブルから`rest_time_seconds`カラムを削除
  - 分析用ビューの更新
- [ ] 既存データの移行テスト

#### 3. **リポジトリ層修正**
- [ ] `strength_repository.go`のSQL修正
  - INSERT文からカテゴリ・休憩時間カラムを削除
  - SELECT文からカテゴリ・休憩時間カラムを削除
- [ ] `strength_query_service.go`のSQL修正
  - PersonalRecordsクエリからカテゴリ・休憩時間を削除
  - findExercises系メソッドからカテゴリ復元処理を削除
  - findSets系メソッドから休憩時間復元処理を削除

#### 4. **アプリケーション層修正**
- [ ] **DTO修正**
  - `ExerciseDTO`から`Category`フィールドを削除
  - `SetDTO`から`RestTimeSeconds`フィールドを削除
- [ ] **マッパー修正**
  - `ToExercise`からカテゴリ変換処理を削除
  - `ToSet`から休憩時間変換処理を削除
  - `FromExercise`からカテゴリ変換処理を削除
  - `FromSet`から休憩時間変換処理を削除
- [ ] **バリデーション修正**
  - カテゴリ・休憩時間のバリデーション処理を削除

#### 5. **MCPツール修正**
- [ ] `training_tool.go`のパラメータ定義修正
  - ツール説明からカテゴリ・休憩時間の記述を削除
  - パラメータ解析処理からカテゴリ・休憩時間を削除
- [ ] レスポンスフォーマッターの修正
  - 応答にカテゴリ・休憩時間が含まれないよう修正

#### 6. **テスト実装**
- [ ] **ドメインテスト**
  - Exercise関連テストの修正・テーブルドリブン化
  - Training関連テストの修正
- [ ] **リポジトリテスト**
  - 保存・取得処理のテスト修正
- [ ] **統合テスト**
  - MCPツール統合テストの修正
  - エンドツーエンドテストの修正

## 影響を受けるファイル一覧

### ドメインレイヤー
- `internal/domain/strength/set.go` ✅
- `internal/domain/strength/set_test.go` ✅
- `internal/domain/strength/exercise.go` ❌
- `internal/domain/strength/exercise_test.go` ❌
- `internal/domain/strength/training_test.go` ❌

### インフラストラクチャレイヤー
- `internal/infrastructure/repository/sqlite/migrations/002_remove_resttime_category.sql` ❌ (新規作成)
- `internal/infrastructure/repository/sqlite/strength_repository.go` ❌
- `internal/infrastructure/query/sqlite/strength_query_service.go` ❌

### アプリケーションレイヤー
- `internal/application/command/dto/strength_command.go` ❌
- `internal/application/command/dto/strength_mapper.go` ❌
- `internal/application/query/dto/training_query.go` ❌
- `internal/application/query/dto/personal_records_dto.go` ❌

### インターフェースレイヤー
- `internal/interface/mcp-tool/tool/training_tool.go` ❌
- `internal/interface/mcp-tool/converter/response_formatter.go` ❌

## 実装優先順位

### 高優先度
1. **Exercise構造体からExerciseCategory削除** - ドメインモデルの整合性確保
2. **データベースマイグレーション作成** - スキーマ変更の準備
3. **リポジトリ層修正** - データアクセス層の修正

### 中優先度
4. **アプリケーション層修正** - DTO・マッパーの修正
5. **MCPツール修正** - ユーザーインターフェースの修正

### 低優先度
6. **テスト実装・実行** - 全体的な動作確認

## 注意事項

### データベース移行
- 既存データが存在する場合、カラム削除前にデータのバックアップが必要
- マイグレーションは可逆的に設計する（ロールバック可能）

### 後方互換性
- この変更はBreaking Changeであり、既存のMCPクライアントとの互換性が失われる
- 必要に応じて移行期間やバージョニング戦略を検討

### テスト戦略
- 各レイヤーの修正完了後に統合テストを実行
- 既存データでのマイグレーションテストを実施
- MCPプロトコル準拠性の確認

## 実装手順

1. ドメインモデル修正（Exercise）
2. データベースマイグレーション作成・テスト
3. リポジトリ層修正
4. アプリケーション層修正
5. MCPツール修正
6. 統合テスト実行・修正
7. 全体動作確認