#\!/bin/bash

# MCPサーバーのテストスクリプト（タイムアウト付き）
# 使用法: ./local_test_query.sh

set -e

echo "🏋️ Fitness MCP Server テスト開始"
echo "=================================="

# サーバーのパス
SERVER_PATH="./mcp"
TIMEOUT=10

# サーバーが存在するかチェック
if [ \! -f "$SERVER_PATH" ]; then
    echo "❌ エラー: $SERVER_PATH が見つかりません"
    echo "まず 'go build -o mcp-server ./cmd/mcp/' を実行してください"
    exit 1
fi

echo "✅ MCPサーバーファイル確認OK"

# テスト関数
test_mcp_call() {
    local test_name="$1"
    local json_request="$2"
    
    echo ""
    echo "🧪 テスト: $test_name"
    echo "リクエスト: $json_request"
    echo ""
    
    # タイムアウト付きでサーバーを実行
    echo "$json_request" | timeout $TIMEOUT "$SERVER_PATH" 2>/dev/null || {
        echo "❌ タイムアウトまたはエラーが発生しました"
        return 1
    }
    
    echo ""
    echo "✅ テスト完了: $test_name"
    echo "----------------------------------------"
}

# 1. ツール一覧の取得
test_mcp_call "ツール一覧取得" '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}'

# 2. 個人記録取得（全て）
test_mcp_call "個人記録取得（全て）" '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_personal_records", "arguments": {}}}'

# 3. 個人記録取得（特定エクササイズ）
test_mcp_call "個人記録取得（ベンチプレス）" '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "get_personal_records", "arguments": {"exercise_name": "ベンチプレス"}}}'

# 4. 期間別トレーニング取得
test_mcp_call "期間別トレーニング取得" '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "get_trainings_by_date_range", "arguments": {"start_date": "2025-06-01", "end_date": "2025-06-30"}}}'

# 5. 存在しないエクササイズでの個人記録取得
test_mcp_call "存在しないエクササイズ" '{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "get_personal_records", "arguments": {"exercise_name": "存在しないエクササイズ"}}}'

echo ""
echo "🎉 全テスト完了！"
echo "=================================="
echo "MCPサーバーが正常に動作しています"
EOF < /dev/null