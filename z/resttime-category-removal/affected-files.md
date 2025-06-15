# 影響を受けるファイル詳細

## ファイル別影響度と修正内容

### 🔴 **高影響度 - 即座に修正が必要**

#### **ドメインレイヤー**

##### `internal/domain/strength/exercise.go`
**影響度:** 🔴 高  
**ステータス:** ❌ 未修正  
**修正内容:**
- `ExerciseCategory`型の削除（15-17行、32-35行、65-89行）
- `Exercise`構造体から`category`フィールド削除（22行）
- `NewExercise`関数の引数変更（92行）
- `Exercise.Category()`メソッド削除（105-107行）
- `Exercise.String()`メソッドからカテゴリ表示削除（154行）

##### `internal/domain/strength/exercise_test.go`
**影響度:** 🔴 高  
**ステータス:** ❌ 未修正  
**修正内容:**
- ExerciseCategory関連テストの削除
- `NewExercise`呼び出しの引数修正
- テーブルドリブンテスト化

##### `internal/domain/strength/training_test.go`
**影響度:** 🔴 高  
**ステータス:** ❌ 未修正  
**修正内容:**
- `NewExercise`呼び出しの引数修正

#### **インフラストラクチャレイヤー**

##### `internal/infrastructure/repository/sqlite/strength_repository.go`
**影響度:** 🔴 高  
**ステータス:** ❌ 未修正  
**修正内容:**
- `saveExercise`メソッド: categoryカラムの削除（200行周辺）
- `saveSet`メソッド: rest_time_secondsカラムの削除（224行周辺）
- INSERT文の修正

##### `internal/infrastructure/query/sqlite/strength_query_service.go`
**影響度:** 🔴 高  
**ステータス:** ❌ 未修正  
**修正内容:**
- `GetPersonalRecords`: カテゴリ・休憩時間取得の削除（198行、212行、227行等）
- `findExercisesByTrainingID`: カテゴリ復元処理削除（413-416行）
- `findSetsByExerciseID`: 休憩時間復元処理削除（558-561行）
- SELECT文の修正

### 🟡 **中影響度 - 段階的修正が必要**

#### **アプリケーションレイヤー**

##### `internal/application/command/dto/strength_command.go`
**影響度:** 🟡 中  
**ステータス:** ❌ 未修正  
**修正内容:**
- `ExerciseDTO.Category`フィールド削除（35行）
- `SetDTO.RestTimeSeconds`フィールド削除（43行）
- バリデーション処理の修正（95-96行、117-118行）

##### `internal/application/command/dto/strength_mapper.go`
**影響度:** 🟡 中  
**ステータス:** ❌ 未修正  
**修正内容:**
- `ToExercise`: カテゴリ変換処理削除（75-79行）
- `ToSet`: 休憩時間変換処理削除（114-117行）
- `FromExercise`: カテゴリ変換処理削除（162行）
- `FromSet`: 休憩時間変換処理削除（178行）

#### **クエリDTO**

##### `internal/application/query/dto/training_query.go`
**影響度:** 🟡 中  
**ステータス:** ❌ 未修正  
**修正内容:**
- レスポンスDTOからカテゴリ・休憩時間フィールド削除

##### `internal/application/query/dto/personal_records_dto.go`
**影響度:** 🟡 中  
**ステータス:** ❌ 未修正  
**修正内容:**
- PersonalRecordsレスポンスからカテゴリ関連情報削除

### 🟢 **低影響度 - 最終段階で修正**

#### **インターフェースレイヤー**

##### `internal/interface/mcp-tool/tool/training_tool.go`
**影響度:** 🟢 低  
**ステータス:** ❌ 未修正  
**修正内容:**
- ツール説明からカテゴリ・休憩時間記述削除（47-54行、60行）
- パラメータ解析処理削除（165-168行、217-221行）

##### `internal/interface/mcp-tool/converter/response_formatter.go`
**影響度:** 🟢 低  
**ステータス:** ❌ 未修正  
**修正内容:**
- フォーマット処理からカテゴリ・休憩時間削除

#### **データベース**

##### `internal/infrastructure/repository/sqlite/migrations/002_remove_resttime_category.sql`
**影響度:** 🟢 低  
**ステータス:** ❌ 新規作成必要  
**修正内容:**
- `exercises.category`カラム削除
- `sets.rest_time_seconds`カラム削除
- 分析用ビューの更新

## 修正完了済みファイル

### ✅ **修正完了**

##### `internal/domain/strength/set.go`
**ステータス:** ✅ 修正完了  
**修正内容:**
- `RestTime`型と関連メソッド削除
- `Set`構造体から`restTime`フィールド削除
- `NewSet`関数引数変更
- `Set.String()`メソッド修正

##### `internal/domain/strength/set_test.go`
**ステータス:** ✅ 修正完了  
**修正内容:**
- テーブルドリブンテスト化
- Given/When/Then形式に変更
- RPEありとRPEなしの両パターンテスト

## 依存関係マップ

```
domain/strength/set.go (✅完了)
    ↓ 依存
domain/strength/exercise.go (❌要修正)
    ↓ 依存
domain/strength/training.go (❌要修正)
    ↓ 依存
application/command/dto/ (❌要修正)
    ↓ 依存
infrastructure/repository/ (❌要修正)
    ↓ 依存
interface/mcp-tool/ (❌要修正)
```

## 修正順序の推奨

### Phase 1: ドメインモデル修正
1. `internal/domain/strength/exercise.go`
2. `internal/domain/strength/exercise_test.go`
3. `internal/domain/strength/training_test.go`

### Phase 2: データベース修正
4. `002_remove_resttime_category.sql`作成

### Phase 3: リポジトリ修正
5. `internal/infrastructure/repository/sqlite/strength_repository.go`
6. `internal/infrastructure/query/sqlite/strength_query_service.go`

### Phase 4: アプリケーション修正
7. `internal/application/command/dto/strength_command.go`
8. `internal/application/command/dto/strength_mapper.go`
9. `internal/application/query/dto/` 配下

### Phase 5: インターフェース修正
10. `internal/interface/mcp-tool/tool/training_tool.go`
11. `internal/interface/mcp-tool/converter/response_formatter.go`

この順序により、依存関係の逆順で修正することで、コンパイルエラーを最小限に抑えて実装できます。