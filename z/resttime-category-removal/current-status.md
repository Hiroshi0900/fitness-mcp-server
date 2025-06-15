# 現在の実装状況と次のステップ

## 現在の状況

### ✅ 完了した作業

#### **Set構造体からRestTime削除 (完了)**

**修正内容:**
- `RestTime`型と関連メソッドを完全削除
- `Set`構造体から`restTime`フィールドを削除
- `NewSet`関数の引数を変更: `(weight, reps, restTime, rpe)` → `(weight, reps, rpe)`
- `Set.String()`メソッドから休憩時間表示を削除

**テスト修正:**
- テーブルドリブンテスト形式に変更
- Given/When/Then構造に変更
- RPEありとRPEなしの両パターンをテスト

**テスト結果:**
```
=== RUN   TestSet_NewSet
=== RUN   TestSet_NewSet/RPEありのセット作成
=== RUN   TestSet_NewSet/RPEなしのセット作成
--- PASS: TestSet_NewSet (0.00s)

=== RUN   TestSet_String
=== RUN   TestSet_String/RPEありのセット文字列表示
=== RUN   TestSet_String/RPEなしのセット文字列表示
--- PASS: TestSet_String (0.00s)
```

**影響範囲:**
- ✅ `internal/domain/strength/set.go`
- ✅ `internal/domain/strength/set_test.go`
- ✅ `internal/domain/strength/exercise_test.go` (依存修正)
- ✅ `internal/domain/strength/training_test.go` (依存修正)

## 次に実装すべき項目

#### **Exercise構造体からExerciseCategory削除 (完了)**

**修正内容:**
- `ExerciseCategory`型と関連メソッドを完全削除
- `Exercise`構造体から`category`フィールドを削除
- `NewExercise`関数の引数を変更: `(name, category)` → `(name)`
- `Exercise.String()`メソッドからカテゴリ表示を削除

**テスト修正:**
- テーブルドリブンテスト形式に変更
- Given/When/Then構造に変更
- 新規テストケース追加（`TestExerciseName_NewExerciseName`, `TestExerciseName_Equals`等）
- ExerciseCategory関連テストを削除

**テスト結果:**
```
PASS
ok      fitness-mcp-server/internal/domain/strength     0.260s
```

**影響範囲:**
- ✅ `internal/domain/strength/exercise.go`
- ✅ `internal/domain/strength/exercise_test.go`
- ✅ `internal/domain/strength/training_test.go` (依存修正)

### 🎯 **優先度2: リポジトリ層修正 (完了)**

**修正内容:**
- `strength_repository.go`のSQL修正
  - `saveExercise`: `category`カラムをINSERT文から削除
  - `saveSet`: `rest_time_seconds`カラムをINSERT文から削除
- `strength_query_service.go`のクエリ修正
  - `findExercisesByTrainingID`: カテゴリ取得・復元処理削除
  - `findExercisesByTrainingIDs`: カテゴリ取得・復元処理削除
  - `findSetsByExerciseID`: 休憩時間取得・復元処理削除
  - `findSetsByExerciseIDs`: 休憩時間取得・復元処理削除
  - `GetPersonalRecords`: カテゴリ・休憩時間関連SQL処理削除

**コンパイル結果:**
- ✅ `internal/infrastructure/repository/sqlite/`
- ✅ `internal/infrastructure/query/sqlite/` (DTO依存エラーは次Phaseで解決)

**影響範囲:**
- ✅ `internal/infrastructure/repository/sqlite/strength_repository.go`
- ✅ `internal/infrastructure/query/sqlite/strength_query_service.go`

## 現在のコンパイル状況

### ✅ 正常
- `internal/domain/strength` パッケージ - 全テストPASS

### ❌ コンパイルエラー予想箇所
以下のパッケージはSet/Exerciseの変更により現在コンパイルエラーが発生している可能性:

1. **アプリケーション層**
   - `internal/application/command/dto/strength_mapper.go`
   - `internal/application/command/handler/strength_handler.go`

2. **インフラストラクチャ層**
   - `internal/infrastructure/repository/sqlite/strength_repository.go`
   - `internal/infrastructure/query/sqlite/strength_query_service.go`

3. **インターフェース層**
   - `internal/interface/mcp-tool/tool/training_tool.go`

## 実装方針

### TDDアプローチ
1. 修正対象の新しいテストケースを先に作成
2. テストが失敗することを確認
3. 実装を修正してテストを通す
4. リファクタリングとテスト追加

### 段階的実装
1. **ドメイン層** → **データベース** → **リポジトリ層** → **アプリケーション層** → **インターフェース層**の順で修正
2. 各段階でコンパイルエラーを解消
3. 統合テストで全体動作を確認

## 推奨次回作業

1. Exercise構造体からExerciseCategory削除のテストケース作成
2. Exercise関連の実装修正
3. データベースマイグレーション作成
4. コンパイルエラー確認と修正方針決定

これにより、ドメインモデルの整合性を保ちながら段階的に修正を進められます。