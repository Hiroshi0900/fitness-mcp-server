# Fitness MCP Server

筋トレ（BIG3中心）とランニング記録管理のためのMCPサーバーです。

## 🎯 プロジェクト目標

- **筋トレ**: ベンチプレス100kg到達（現在95kg、体重72kg）
- **ランニング**: ハーフマラソン出場（7/12）、将来的にフルマラソン
- **完成期限**: 2025年8月まで

## 🛠️ 技術スタック

- **言語**: Go 1.23.4
- **アーキテクチャ**: クリーンアーキテクチャ
- **データベース**: SQLite（ローカル、ファイルベース）
- **MCPライブラリ**: [mcp-go](https://github.com/mark3labs/mcp-go)
- **開発手法**: TDD + 軽量DDD

## 📁 プロジェクト構造

```
fitness-mcp-server/
├── cmd/mcp/              # MCPサーバーのエントリーポイント
├── internal/
│   ├── application/      # アプリケーション層
│   │   ├── command/      # コマンド（書き込み操作）
│   │   └── query/        # クエリ（読み込み操作）
│   ├── domain/           # ドメイン層
│   │   ├── strength/     # 筋トレドメイン
│   │   └── running/      # ランニングドメイン
│   ├── infrastructure/   # インフラ層
│   └── interface/        # インターフェース層
├── data/                 # SQLiteデータベースファイル
├── docker-compose.yml    # Docker Compose設定
├── Dockerfile           # Docker設定
└── Makefile            # 開発用コマンド
```

## 🚀 クイックスタート

### 前提条件

- Docker & Docker Compose
- Go 1.23+ （ローカル開発時）
- Make （オプション、コマンド簡略化のため）

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd fitness-mcp-server
```

### 2. 開発環境のセットアップ

```bash
# Makefileを使用する場合
make dev-setup

# 直接実行する場合
./docker-run.sh build
```

### 3. MCPサーバーの起動

```bash
# Dockerで起動
make run
# または
./docker-run.sh run

# ローカルで起動（要Go環境）
make build
./mcp
```

### 4. 動作確認

```bash
# 全テスト実行
make docker-test

# 記録テスト
make test-record

# クエリテスト  
make test-query
```

## 📋 使用可能なMCPツール

MCPサーバーは以下のツールを提供します：

### 1. record_training - トレーニング記録

筋トレセッションを記録します。

```json
{
  \"jsonrpc\": \"2.0\",
  \"id\": 1,
  \"method\": \"tools/call\",
  \"params\": {
    \"name\": \"record_training\",
    \"arguments\": {
      \"date\": \"2025-06-14\",
      \"exercises\": [
        {
          \"name\": \"ベンチプレス\",
          \"category\": \"Compound\",
          \"sets\": [
            {
              \"weight_kg\": 95,
              \"reps\": 8,
              \"rest_time_seconds\": 180,
              \"rpe\": 8
            }
          ]
        }
      ],
      \"notes\": \"調子良好\"
    }
  }
}
```

### 2. get_personal_records - 個人記録取得

個人記録（PR）を取得します。

```json
{
  \"jsonrpc\": \"2.0\",
  \"id\": 2,
  \"method\": \"tools/call\",
  \"params\": {
    \"name\": \"get_personal_records\",
    \"arguments\": {
      \"exercise_name\": \"ベンチプレス\"  // オプション、指定しない場合は全エクササイズ
    }
  }
}
```

### 3. get_trainings_by_date_range - 期間別トレーニング取得

指定期間のトレーニング履歴を取得します。

```json
{
  \"jsonrpc\": \"2.0\",
  \"id\": 3,
  \"method\": \"tools/call\",
  \"params\": {
    \"name\": \"get_trainings_by_date_range\",
    \"arguments\": {
      \"start_date\": \"2025-06-01\",
      \"end_date\": \"2025-06-30\"
    }
  }
}
```

## 🔧 MCPクライアント接続設定

### 自動設定（推奨）

**ローカルバイナリを使用する場合:**
```bash
# ローカルでビルド & Claude Desktop設定
make setup-claude-local
```

**Dockerを使用する場合:**
```bash  
# Dockerイメージビルド & Claude Desktop設定
make setup-claude-docker
```

**設定確認:**
```bash
make check-claude-config
```

### 手動設定

#### Claude Desktop での設定

Claude Desktopの設定ファイル（`config.json`）に以下を追加：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\\Claude\\claude_desktop_config.json`

```json
{
  \"mcpServers\": {
    \"fitness-mcp-server\": {
      \"command\": \"docker\",
      \"args\": [
        \"run\",
        \"--rm\",
        \"-i\",
        \"-v\",
        \"/path/to/fitness-mcp-server/data:/app/data\",
        \"fitness-mcp-server_fitness-mcp\"
      ],
      \"env\": {
        \"MCP_DATA_DIR\": \"/app/data\"
      }
    }
  }
}
```

**注意**: `/path/to/fitness-mcp-server`は実際のプロジェクトパスに置き換えてください。

### VS Code MCP Extension での設定

VS CodeのMCP拡張機能を使用する場合：

```json
{
  \"mcp.servers\": {
    \"fitness-mcp-server\": {
      \"command\": \"./mcp\",
      \"args\": [],
      \"cwd\": \"/path/to/fitness-mcp-server\"
    }
  }
}
```

## 💻 開発用コマンド

### Makefileコマンド一覧

```bash
# ヘルプ表示
make help

# 開発環境
make dev-setup      # 開発環境セットアップ
make docker-dev     # 開発モード（シェル）起動
make build          # ローカルビルド
make clean          # クリーンアップ

# サーバー操作
make run            # サーバー起動
make stop           # サーバー停止
make restart        # サーバー再起動
make logs           # ログ表示
make status         # 状態確認

# テスト
make test           # Goユニットテスト
make docker-test    # Docker環境でMCPテスト
make test-local     # ローカルでMCPテスト
make test-record    # 記録テストのみ
make test-query     # クエリテストのみ

# 品質チェック
make lint           # 静的解析
make check          # lint + test
make ci             # CI用（全チェック）

# その他
make info           # プロジェクト情報表示
make docker-clean   # Dockerリソース削除
```

### 直接実行コマンド

```bash
# Docker操作
./docker-run.sh build     # イメージビルド
./docker-run.sh run       # サーバー起動
./docker-run.sh dev       # 開発モード
./docker-run.sh test      # テスト実行
./docker-run.sh clean     # クリーンアップ

# Docker テスト
./docker-test.sh test-all       # 全テスト
./docker-test.sh test-record    # 記録テスト
./docker-test.sh test-query     # クエリテスト
./docker-test.sh interactive    # インタラクティブモード

# ローカルテスト
./local_test_query.sh     # ローカルMCPテスト
./test_commands.sh        # 基本的なコマンドテスト
```

## 🏗️ アーキテクチャ

### 設計原則

- **クリーンアーキテクチャ**: 依存関係の逆転により保守性を向上
- **CQRS軽量適用**: コマンド（書き込み）とクエリ（読み込み）の分離
- **軽量DDD**: Entity、Value Object、Repository パターンを採用
- **TDD**: テスト駆動開発による品質保証

### レイヤー構成

1. **Interface Layer**: MCPプロトコル、外部API
2. **Application Layer**: ユースケース、DTO、ハンドラー
3. **Domain Layer**: ビジネスロジック、エンティティ
4. **Infrastructure Layer**: データベース、外部サービス

## 🗄️ データベース

### SQLite選択理由

- **コスト**: 完全無料
- **シンプルさ**: ファイルベース、サーバー不要
- **ポータビリティ**: 単一ファイルでのデータ管理
- **Go対応**: 優秀なドライバー（modernc.org/sqlite）

### マイグレーション

データベーススキーマは `internal/infrastructure/repository/sqlite/migrations/` に配置。

## 🧪 テスト

### テスト種類

- **ユニットテスト**: `go test ./...`
- **MCPテスト**: JSON-RPCプロトコルでの動作確認
- **統合テスト**: Docker環境での全体動作確認

### テスト実行方法

```bash
# ユニットテスト
make test

# MCPプロトコルテスト
make docker-test

# 全テスト
make ci
```

## 🚨 トラブルシューティング

### よくある問題

1. **Dockerイメージがビルドできない**
   ```bash
   make docker-clean
   make docker-build
   ```

2. **データベースファイルが見つからない**
   ```bash
   # dataディレクトリが存在することを確認
   ls -la ./data/
   ```

3. **MCPサーバーに接続できない**
   - Docker Composeが正常に起動していることを確認
   - ポートの競合がないことを確認
   - ログを確認: `make logs`

### デバッグ方法

```bash
# 詳細ログ表示
make logs

# 開発モードでシェル起動
make docker-dev

# サーバー状態確認
make status
```

## 📝 ライセンス

このプロジェクトは個人利用目的で開発されています。

## 🤝 コントリビューション

現在は個人プロジェクトのため、外部からのコントリビューションは受け付けていません。

## 📧 連絡先

プロジェクトに関する質問や提案がある場合は、GitHubのIssueを作成してください。
