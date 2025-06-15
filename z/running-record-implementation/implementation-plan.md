# ランニング記録機能実装計画

## 1. 設計方針

### アーキテクチャ判断
**別テーブル設計を採用する理由:**
- 筋トレとランニングは根本的に異なるドメイン
- データ構造が全く違う（重量・回数 vs 距離・時間・ペース）
- 将来の拡張性を考慮（心拍数ゾーン分析、ルート情報等）
- SQLの最適化が容易

### ドメイン分離
- `strength`パッケージ：筋トレ専用
- `running`パッケージ：ランニング専用
- `shared`パッケージ：共通ID等

## 2. データベース設計

### ランニング専用テーブル
```sql
-- ランニングセッションテーブル
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

-- インデックス
CREATE INDEX IF NOT EXISTS idx_running_sessions_date ON running_sessions(date);
CREATE INDEX IF NOT EXISTS idx_running_sessions_run_type ON running_sessions(run_type);
CREATE INDEX IF NOT EXISTS idx_running_sessions_distance ON running_sessions(distance_km);

-- 分析用ビュー
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

## 3. 実装手順（後方互換性重視）

### Phase 1: 基盤整備（既存機能への影響なし）
- [ ] `003_add_running_tables.sql`マイグレーション作成
- [ ] マイグレーション実行・検証
- [ ] 既存DBとの共存確認

### Phase 2: Repository層（独立実装）
- [ ] `RunningRepository`インターフェース定義
- [ ] SQLite実装クラス作成（新テーブルのみ使用）
- [ ] 既存Repositoryとの分離確認

### Phase 3: Application層（独立実装）
- [ ] `RecordRunningCommand`DTO作成
- [ ] `RunningCommandHandler`実装
- [ ] 既存ハンドラーとの分離確認

### Phase 4: MCPツール追加（既存ツールと並行）
- [ ] `record_running`ツール作成
- [ ] 既存`record_training`ツールとの並行動作確認
- [ ] エラーハンドリング独立性確認

### Phase 5: 統合検証
- [ ] 両機能の並行動作テスト
- [ ] データベースファイル整合性確認
- [ ] 既存機能の回帰テスト

## 4. MCPツール仕様設計

### `record_running`ツール
```json
{
  "name": "record_running",
  "description": "ランニングセッションを記録するツール",
  "parameters": {
    "date": "実施日（YYYY-MM-DD）",
    "distance_km": "距離（km）",
    "duration": "時間（MM:SS形式または分数）",
    "run_type": "Easy|Tempo|Interval|Long|Race",
    "heart_rate_bpm": "心拍数（オプション）",
    "notes": "メモ（オプション）"
  }
}
```

### 使用例
```bash
# MM:SS形式での入力
record_running \
  --date 2025-06-16 \
  --distance_km 5.0 \
  --duration "25:30" \
  --run_type Easy \
  --notes "気持ちよく走れた"

# 分数での入力も対応
record_running \
  --date 2025-06-16 \
  --distance_km 8.0 \
  --duration 35.5 \
  --run_type Tempo \
  --heart_rate_bpm 165
```

## 5. エラーハンドリング

### 入力フォーマット対応
- **時間入力**: MM:SS形式（例：25:30）または分数（例：25.5）の両方に対応
- **自動パース**: 文字列に`:`が含まれていればMM:SS、そうでなければ分数として処理

### 基本的なバリデーション
- 距離: 正の数値
- 時間: 正の値
- ランニングタイプ: 定義済み値のみ
- 心拍数: 正の整数（オプション）

### エラーメッセージ例
```
❌ 時間の形式が不正です: "25:ab"（MM:SS形式または分数で入力してください）
❌ 無効なランニングタイプ: "Jog"（有効値: Easy, Tempo, Interval, Long, Race）
❌ 心拍数は正の整数で入力してください: "-10"
```

## 6. 技術的考慮事項

### パフォーマンス
- 日付・タイプ・距離にインデックス作成
- 分析用ビューで集計クエリ最適化

### 拡張性
- 将来のGPS軌跡データに備えたテーブル設計
- 天候・路面状況などの環境要因追加余地

### データ整合性
- 外部キー制約でデータ整合性保証
- トランザクション処理でatomicity確保

## 7. 品質保証

### テスト項目
- [ ] ドメインモデルの単体テスト
- [ ] Repository層の統合テスト
- [ ] MCPツールのE2Eテスト
- [ ] バリデーション境界値テスト

### レビューポイント
- [ ] ドメインロジックの適切な配置
- [ ] エラーハンドリングの網羅性
- [ ] SQLインジェクション対策
- [ ] 日本語メッセージの品質

## 8. 運用考慮

### ログ・監視
- 記録成功/失敗のログ出力
- パフォーマンスメトリクス収集

### バックアップ
- ランニングデータの定期バックアップ
- 既存筋トレデータとの整合性維持

---

この計画に基づいて、段階的に実装を進めることで、既存の筋トレ機能に影響を与えることなく、ランニング記録機能を追加できます。
