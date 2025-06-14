# Fitness MCP Server Makefile
# 筋トレ・ランニング記録管理MCPサーバー開発用

.PHONY: help build test clean dev run stop logs status lint check docker-build docker-test docker-dev

# デフォルトターゲット
help: ## ヘルプを表示
	@echo "Fitness MCP Server - 開発用Makefile"
	@echo ""
	@echo "使用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# =========================================
# ローカル開発用
# =========================================

build: ## ローカルでバイナリをビルド
	@echo "🔨 バイナリをビルド中..."
	go build -o mcp ./cmd/mcp/
	@echo "✅ ビルド完了: ./mcp"

test: ## ローカルでGoテストを実行
	@echo "🧪 Goテスト実行中..."
	go test ./...
	@echo "✅ テスト完了"

lint: ## コードの静的解析
	@echo "🔍 静的解析実行中..."
	go vet ./...
	go fmt ./...
	@echo "✅ 静的解析完了"

check: lint test ## lintとtestを実行（コンパイルチェック含む）

clean: ## ビルド成果物とキャッシュを削除
	@echo "🧹 クリーンアップ中..."
	rm -f mcp
	go clean -cache
	@echo "✅ クリーンアップ完了"

# =========================================
# Docker開発用
# =========================================

docker-build: ## Dockerイメージをビルド
	@echo "🐳 Dockerイメージビルド中..."
	./docker-run.sh build

docker-dev: ## 開発モード（シェル）でDockerコンテナを起動
	@echo "🚀 開発モード起動中..."
	./docker-run.sh dev

run: ## MCPサーバーをDockerで起動
	@echo "🏃 MCPサーバー起動中..."
	./docker-run.sh run

stop: ## MCPサーバーを停止
	@echo "⏹️  MCPサーバー停止中..."
	./docker-run.sh stop

restart: ## MCPサーバーを再起動
	@echo "🔄 MCPサーバー再起動中..."
	./docker-run.sh restart

logs: ## MCPサーバーのログを表示
	./docker-run.sh logs

status: ## 現在の状態を確認
	./docker-run.sh status

# =========================================
# テスト用
# =========================================

docker-test: ## Dockerでテスト実行（記録+クエリ）
	@echo "🧪 Docker環境でテスト実行中..."
	./docker-test.sh test-all

test-record: ## 記録テストのみ実行
	./docker-test.sh test-record

test-query: ## クエリテストのみ実行
	./docker-test.sh test-query

test-local: build ## ローカルでMCPテストを実行
	@echo "🧪 ローカルMCPテスト実行中..."
	./local_test_query.sh

test-interactive: ## インタラクティブテストモード
	./docker-test.sh interactive

# =========================================
# 総合コマンド
# =========================================

dev-setup: docker-build ## 開発環境をセットアップ
	@echo "🛠️  開発環境セットアップ完了"
	@echo "次のコマンドで開発を開始できます:"
	@echo "  make docker-dev  # 開発モード"
	@echo "  make run         # サーバー起動"
	@echo "  make docker-test # テスト実行"

ci: check test-local ## CI/CD用（lint + test + ローカルMCPテスト）

docker-clean: ## Dockerリソースを完全削除
	./docker-run.sh clean

# =========================================
# Claude Desktop設定
# =========================================

setup-claude-local: build ## Claude Desktop用ローカル設定
	./setup-claude.sh setup-local

setup-claude-docker: docker-build ## Claude Desktop用Docker設定  
	./setup-claude.sh setup-docker

check-claude-config: ## Claude Desktop設定を確認
	./setup-claude.sh check

# =========================================
# 情報表示
# =========================================

info: ## プロジェクト情報を表示
	@echo "📊 Fitness MCP Server 情報"
	@echo "================================"
	@echo "Go Version: $$(go version)"
	@echo "Module: $$(head -1 go.mod)"
	@echo "Docker: $$(docker --version 2>/dev/null || echo 'Not installed')"
	@echo "Docker Compose: $$(docker-compose --version 2>/dev/null || echo 'Not installed')"
	@echo ""
	@echo "📁 プロジェクト構造:"
	@find . -name "*.go" -o -name "*.sql" -o -name "Dockerfile" -o -name "docker-compose.yml" | head -10
	@echo ""
	@echo "🎯 目標:"
	@echo "  - ベンチプレス: 100kg到達（現在95kg）"
	@echo "  - ハーフマラソン: 7/12出場"
