#!/bin/bash

echo "=== MCP Fitness Server Test ==="

# 初期化
echo "1. 初期化中..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./mcp

# ツール一覧
echo -e "\n2. 利用可能なツール一覧..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./mcp

# クエリテスト
echo -e "\n3. 期間検索テスト（正常ケース）..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_trainings_by_date_range","arguments":{"start_date":"2025-06-01","end_date":"2025-06-30"}}}' | ./mcp

echo -e "\n4. 期間制限テスト（エラーケース）..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_trainings_by_date_range","arguments":{"start_date":"2024-01-01","end_date":"2025-12-31"}}}' | ./mcp

echo -e "\n=== テスト完了 ==="
