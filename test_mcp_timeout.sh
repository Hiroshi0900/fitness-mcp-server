#!/bin/bash

# タイムアウト付きMCPテストスクリプト
TIMEOUT_SECONDS=10

echo "=== MCP Server Test with Timeout ==="

# 1. ビルド
echo "1. Building MCP server..."
go build -o mcp cmd/mcp/main.go
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"

# 2. タイムアウト付きクエリテスト
echo ""
echo "2. Testing get_trainings_by_date_range with timeout..."

# JSON要求を作成
REQUEST='{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_trainings_by_date_range","arguments":{"start_date":"2025-06-01","end_date":"2025-06-30"}}}'

echo "Request: $REQUEST"
echo ""
echo "Testing with ${TIMEOUT_SECONDS}s timeout..."

# バックグラウンドでMCPサーバーを実行し、PIDを取得
echo "$REQUEST" | gtimeout ${TIMEOUT_SECONDS}s ./mcp

# タイムアウトの場合は124、通常終了は0
EXIT_CODE=$?

echo ""
if [ $EXIT_CODE -eq 124 ]; then
    echo "⏰ Request timed out after ${TIMEOUT_SECONDS} seconds"
elif [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Request completed successfully"
else
    echo "❌ Request failed with exit code: $EXIT_CODE"
fi

echo ""
echo "3. Testing shorter timeout (5s)..."
echo "$REQUEST" | gtimeout 5s ./mcp

EXIT_CODE=$?
if [ $EXIT_CODE -eq 124 ]; then
    echo "⏰ Request timed out after 5 seconds"
elif [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Request completed successfully"
else
    echo "❌ Request failed with exit code: $EXIT_CODE"
fi

echo ""
echo "=== Test Complete ===" 
