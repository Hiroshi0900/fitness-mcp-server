# ランニング記録機能実装プロジェクト

## 概要
筋トレ記録システムにランニング記録機能を追加する実装プロジェクト

## 基本方針
- **後方互換性重視**: 既存の筋トレ機能に影響を与えない
- **ドメイン分離**: 筋トレとランニングは完全に独立したテーブル設計
- **段階的実装**: Phase毎に既存機能との共存を確認

## ディレクトリ構成
```
z/running-record-implementation/
├── README.md                 # このファイル
├── implementation-plan.md    # 詳細実装計画
├── task-list.md             # タスク一覧と進捗
├── db-design.md             # データベース設計
├── api-design.md            # MCPツール設計
└── progress-log.md          # 実装進捗ログ
```

## クイックスタート
1. `task-list.md` で現在のタスクを確認
2. `implementation-plan.md` で詳細計画を参照
3. `progress-log.md` で進捗状況を記録

## 関連ファイル
- 既存筋トレドメイン: `internal/domain/strength/`
- 既存ランニングドメイン: `internal/domain/running/`
- マイグレーション: `internal/infrastructure/repository/sqlite/migrations/`
- MCPツール: `internal/interface/mcp-tool/tool/`

## 連絡先
実装中の質問や相談は随時相談してください。
