#!/bin/bash

echo "=== MCP Fitness Server 動作確認 ==="
echo

echo "1. サーバー初期化テスト..."
echo '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}' | ./mcp-server
echo
echo "---"

echo "2. ツール一覧確認..."
echo '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":2}' | ./mcp-server
echo
echo "---"

echo "3. ベンチプレス記録テスト..."
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"record_training","arguments":{"date":"2025-01-10T09:00:00Z","exercises":[{"name":"ベンチプレス","category":"Compound","sets":[{"weight_kg":95.0,"reps":8,"rest_time_seconds":180,"rpe":9}]}],"notes":"100kg目標まであと5kg！"}},"id":3}' | ./mcp-server
echo
echo "---"

echo "4. スクワット記録テスト..."
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"record_training","arguments":{"date":"2025-01-11T10:00:00Z","exercises":[{"name":"スクワット","category":"Compound","sets":[{"weight_kg":80.0,"reps":10,"rest_time_seconds":120}]}],"notes":"フォーム重視で実施"}},"id":4}' | ./mcp-server
echo
echo "---"

echo "✅ 動作確認完了！"
