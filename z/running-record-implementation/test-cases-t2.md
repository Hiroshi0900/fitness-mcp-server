# T2: Repository層実装 - テストケース洗い出し（記録専用）

## 対象
- `RunningRepository`インターフェース定義（Saveのみ）
- SQLite実装クラス作成

## テストレベル
- **ITテスト**: 実際のSQLite接続を使用（Repository層の特性上）
- **ユニットテスト**: ドメインモデル層で実施済み

## テストケース一覧

### T2-1: RunningRepositoryインターフェース

#### インターフェース設計（最小版）
```go
type RunningRepository interface {
    // Save はランニングセッションを保存します
    Save(session *running.RunningSession) error
}
```

### T2-2: SQLite実装ITテストケース

#### TC2-1: Save操作 - 正常系
```go
func TestSQLiteRunningRepository_Save_ValidSession(t *testing.T) {
    // Given: 有効なRunningSessionと初期化済みRepository
    // When: Save実行
    // Then: エラーなし、DB内にデータ存在確認
}

func TestSQLiteRunningRepository_Save_WithHeartRate(t *testing.T) {
    // Given: 心拍数付きRunningSession
    // When: Save実行
    // Then: 心拍数含めて正常保存
}

func TestSQLiteRunningRepository_Save_WithoutHeartRate(t *testing.T) {
    // Given: 心拍数なしRunningSession
    // When: Save実行
    // Then: 心拍数NULLで正常保存
}
```

#### TC2-2: Save操作 - 異常系
```go
func TestSQLiteRunningRepository_Save_DuplicateID(t *testing.T) {
    // Given: 同じIDのセッションが既にDB保存済み
    // When: 同じIDで再度Save実行
    // Then: 主キー制約エラー発生
}

func TestSQLiteRunningRepository_Save_InvalidDistance(t *testing.T) {
    // Given: 距離が0のRunningSession
    // When: Save実行
    // Then: CHECK制約エラー発生
}

func TestSQLiteRunningRepository_Save_InvalidDuration(t *testing.T) {
    // Given: 時間が0のRunningSession
    // When: Save実行
    // Then: CHECK制約エラー発生
}

func TestSQLiteRunningRepository_Save_InvalidRunType(t *testing.T) {
    // Given: 不正なrun_typeのRunningSession
    // When: Save実行
    // Then: CHECK制約エラー発生
}
```

#### TC2-2: Update操作
```go
func TestSQLiteRunningRepository_Update_ExistingSession(t *testing.T) {
    // Arrange: 既存セッション
    // Act: 内容変更してUpdate
    // Assert: 正常更新、変更内容確認
}

func TestSQLiteRunningRepository_Update_NonExistentSession(t *testing.T) {
    // Arrange: 存在しないID
    // Act: Update実行
    // Assert: エラー発生
}
```

#### TC2-3: Delete操作
```go
func TestSQLiteRunningRepository_Delete_ExistingSession(t *testing.T) {
    // Arrange: 既存セッション
    // Act: Delete実行
    // Assert: 削除成功、DB確認
}

func TestSQLiteRunningRepository_Delete_NonExistentSession(t *testing.T) {
    // Arrange: 存在しないID
    // Act: Delete実行
    // Assert: エラー発生
}
```

#### TC2-4: FindByID操作
```go
func TestSQLiteRunningRepository_FindByID_ExistingSession(t *testing.T) {
    // Arrange: 保存済みセッション
    // Act: FindByID実行
    // Assert: 正しいデータ取得
}

func TestSQLiteRunningRepository_FindByID_NonExistentSession(t *testing.T) {
    // Arrange: 存在しないID
    // Act: FindByID実行
    // Assert: エラー発生
}
```

#### TC2-5: データ型変換テスト
```go
func TestSQLiteRunningRepository_DataConversion_Distance(t *testing.T) {
    // Arrange: 距離データを含むセッション
    // Act: Save→FindByID
    // Assert: 距離値が正確に保存・復元
}

func TestSQLiteRunningRepository_DataConversion_Duration(t *testing.T) {
    // Arrange: 時間データを含むセッション
    // Act: Save→FindByID  
    // Assert: 時間値が正確に保存・復元
}

func TestSQLiteRunningRepository_DataConversion_Pace(t *testing.T) {
    // Arrange: ペースデータを含むセッション
    // Act: Save→FindByID
    // Assert: ペース値が正確に保存・復元
}

func TestSQLiteRunningRepository_DataConversion_HeartRate(t *testing.T) {
    // Arrange: 心拍数データ（nil、有効値）
    // Act: Save→FindByID
    // Assert: 心拍数が正確に保存・復元
}

func TestSQLiteRunningRepository_DataConversion_RunType(t *testing.T) {
    // Arrange: 各ランニングタイプ
    // Act: Save→FindByID
    // Assert: タイプが正確に保存・復元
}
```

#### TC2-6: トランザクション処理
```go
func TestSQLiteRunningRepository_Transaction_SaveRollback(t *testing.T) {
    // Arrange: 不正データでSave失敗想定
    // Act: Save実行（失敗）
    // Assert: ロールバック、データベース未変更
}

func TestSQLiteRunningRepository_Transaction_UpdateRollback(t *testing.T) {
    // Arrange: 不正データでUpdate失敗想定
    // Act: Update実行（失敗）
    // Assert: ロールバック、元データ保持
}
```

#### TC2-7: 並行処理・エラーハンドリング
```go
func TestSQLiteRunningRepository_ConcurrentAccess(t *testing.T) {
    // Arrange: 複数goroutineでアクセス
    // Act: 同時Save/Update/Delete
    // Assert: データ整合性保持
}

func TestSQLiteRunningRepository_DatabaseConnectionError(t *testing.T) {
    // Arrange: 無効なDB接続
    // Act: Repository操作
    // Assert: 適切なエラーハンドリング
}
```

### T2-3: 既存Repositoryとの分離確認

#### TC2-8: 独立性確認
```go
func TestRunningRepository_Independence_FromStrength(t *testing.T) {
    // Arrange: 筋トレ・ランニング両Repository
    // Act: 両方で操作実行
    // Assert: 相互に影響しない
}

func TestRunningRepository_SharedConnection_Safety(t *testing.T) {
    // Arrange: 同一DB接続を共有
    // Act: 筋トレ・ランニングの並行操作
    // Assert: データ破損なし
}
```

## テストヘルパー設計

### テスト用ユーティリティ
```go
// テスト用RunningSession作成
func createTestRunningSession(t *testing.T) *running.RunningSession

// テスト用インメモリDB作成
func setupTestDB(t *testing.T) *sql.DB

// テスト後クリーンアップ  
func cleanupTestDB(t *testing.T, db *sql.DB)

// DBの状態確認ヘルパー
func assertSessionExistsInDB(t *testing.T, db *sql.DB, id shared.SessionID)
func assertSessionNotExistsInDB(t *testing.T, db *sql.DB, id shared.SessionID)

// データ比較ヘルパー
func assertSessionEquals(t *testing.T, expected, actual *running.RunningSession)
```

## 実装ファイル計画

### ファイル構成
```
internal/interface/repository/
├── running_repository.go           # インターフェース定義

internal/infrastructure/repository/sqlite/
├── running_repository.go           # SQLite実装
├── running_repository_test.go      # 実装テスト
```

### 実装順序
1. **Red**: テスト実装（全てコメントアウト）
2. **Green**: インターフェース定義 → SQLite最小実装
3. **Refactor**: エラーハンドリング・パフォーマンス改善

## 次のアクション

このテストケース洗い出しで問題なければ：
1. テスト実装（コメントアウト状態）
2. `RunningRepository`インターフェース定義
3. SQLite実装クラス作成（最小限）
4. テスト通過までリファクタリング

問題や追加すべき観点があればお知らせください。
