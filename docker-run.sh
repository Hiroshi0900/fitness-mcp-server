#!/bin/bash

# fitness-mcp-server Docker 実行スクリプト

set -e

# 色付きメッセージ用の関数
log_info() {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}

log_success() {
    echo -e "\033[1;32m[SUCCESS]\033[0m $1"
}

log_error() {
    echo -e "\033[1;31m[ERROR]\033[0m $1"
}

log_warning() {
    echo -e "\033[1;33m[WARNING]\033[0m $1"
}

# 使用方法の表示
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build     - Docker imageをビルドする"
    echo "  run       - MCPサーバーを起動する"
    echo "  stop      - MCPサーバーを停止する"
    echo "  restart   - MCPサーバーを再起動する"
    echo "  logs      - ログを表示する"
    echo "  dev       - 開発モード（シェル）で起動する"
    echo "  test      - MCPサーバーでテストを実行する"
    echo "  clean     - 停止してイメージを削除する"
    echo "  status    - 現在の状態を確認する"
    echo "  help      - このヘルプを表示する"
}

# Docker, Docker Composeの存在確認
check_prerequisites() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker が見つかりません。Dockerをインストールしてください。"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose が見つかりません。Docker Composeをインストールしてください。"
        exit 1
    fi
}

# データディレクトリの確認・作成
ensure_data_directory() {
    if [ ! -d "./data" ]; then
        log_info "データディレクトリを作成しています..."
        mkdir -p ./data
        log_success "データディレクトリを作成しました: ./data"
    fi
}

# Docker imageのビルド
build_image() {
    log_info "Docker imageをビルドしています..."
    docker-compose build fitness-mcp
    log_success "Docker imageのビルドが完了しました"
}

# MCPサーバーの起動
start_server() {
    ensure_data_directory
    log_info "MCPサーバーを起動しています..."
    docker-compose up -d fitness-mcp
    log_success "MCPサーバーが起動しました"
    log_info "ログを確認するには: $0 logs"
}

# MCPサーバーの停止
stop_server() {
    log_info "MCPサーバーを停止しています..."
    docker-compose down
    log_success "MCPサーバーを停止しました"
}

# MCPサーバーの再起動
restart_server() {
    log_info "MCPサーバーを再起動しています..."
    stop_server
    start_server
}

# ログの表示
show_logs() {
    docker-compose logs -f fitness-mcp
}

# 開発モードの起動（シェル）
start_dev() {
    ensure_data_directory
    log_info "開発モード（シェル）でコンテナを起動しています..."
    docker-compose --profile dev run --rm fitness-mcp-dev
}

# テスト実行
run_test() {
    ensure_data_directory
    log_info "MCPサーバーでテストを実行します..."
    
    # テストJSON
    TEST_JSON='{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "record_training", "arguments": {"date": "2025-06-13", "exercises": [{"name": "ベンチプレス", "category": "Compound", "sets": [{"weight_kg": 95, "reps": 8, "rest_time_seconds": 180}]}], "notes": "テスト実行"}}}'
    
    log_info "記録テストを実行中..."
    echo "$TEST_JSON" | docker-compose run --rm -T fitness-mcp
    log_success "テスト完了"
}

# クリーンアップ
clean_up() {
    log_info "停止とクリーンアップを実行しています..."
    docker-compose down
    docker-compose down --rmi local --volumes --remove-orphans
    log_success "クリーンアップが完了しました"
}

# 状態確認
show_status() {
    log_info "現在の状態:"
    echo ""
    echo "=== Docker Containers ==="
    docker-compose ps
    echo ""
    echo "=== Data Directory ==="
    if [ -d "./data" ]; then
        ls -la ./data/
    else
        echo "データディレクトリが存在しません"
    fi
}

# メイン処理
main() {
    check_prerequisites

    case "${1:-help}" in
        build)
            build_image
            ;;
        run)
            start_server
            ;;
        stop)
            stop_server
            ;;
        restart)
            restart_server
            ;;
        logs)
            show_logs
            ;;
        dev)
            start_dev
            ;;
        test)
            run_test
            ;;
        clean)
            clean_up
            ;;
        status)
            show_status
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            log_error "不明なコマンド: $1"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# スクリプト実行
main "$@"
