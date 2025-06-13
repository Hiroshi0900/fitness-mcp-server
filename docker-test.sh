#!/bin/bash

# Docker環境でのMCPサーバーテスト用スクリプト

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

# テストデータ - 1行のJSONとして定義
TEST_JSON='{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "record_training", "arguments": {"date": "2025-06-13", "exercises": [{"name": "ベンチプレス", "category": "Compound", "sets": [{"weight_kg": 95, "reps": 8, "rest_time_seconds": 180, "rpe": 8}]}], "notes": "Docker環境でのテスト"}}}'

# クエリテスト用JSON
QUERY_JSON='{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_trainings_by_date_range", "arguments": {"start_date": "2025-06-01", "end_date": "2025-06-30"}}}'

# 使用方法の表示
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  test-record   - 記録テストを実行"
    echo "  test-query    - クエリテストを実行"
    echo "  test-all      - 全テストを実行"
    echo "  interactive   - インタラクティブモードで起動"
    echo "  help          - このヘルプを表示"
}

# 記録テスト
test_record() {
    log_info "記録テストを実行中..."
    # メインのfitness-mcpサービスを使用
    echo "$TEST_JSON" | docker-compose run --rm -T fitness-mcp
    log_success "記録テスト完了"
}

# クエリテスト
test_query() {
    log_info "クエリテストを実行中..."
    echo "$QUERY_JSON" | docker-compose run --rm -T fitness-mcp
    log_success "クエリテスト完了"
}

# 全テスト実行
test_all() {
    log_info "全テストを開始..."
    test_record
    echo ""
    test_query
    log_success "全テスト完了"
}

# インタラクティブモード
interactive_mode() {
    log_info "インタラクティブモードで起動中..."
    log_info "JSON-RPCリクエストを入力してください（Ctrl+Dで終了）"
    docker-compose run --rm -T fitness-mcp
}

# メイン処理
main() {
    # データディレクトリの確認
    if [ ! -d "./data" ]; then
        log_info "データディレクトリを作成しています..."
        mkdir -p ./data
    fi

    case "${1:-help}" in
        test-record)
            test_record
            ;;
        test-query)
            test_query
            ;;
        test-all)
            test_all
            ;;
        interactive)
            interactive_mode
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
