#!/bin/bash

echo "=== curl でMCP APIテスト ==="
echo

echo "1. ヘルスチェック"
curl -s http://localhost:8080/health | jq '.'
echo
echo "---"

echo "2. 初期化テスト"
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}' | jq '.'
echo
echo "---"

echo "3. ツール一覧確認"
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":2}' | jq '.'
echo
echo "---"

echo "4. ベンチプレス記録"
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"tools/call",
    "params":{
      "name":"record_training",
      "arguments":{
        "date":"2025-01-12T10:00:00Z",
        "exercises":[{
          "name":"ベンチプレス",
          "category":"Compound",
          "sets":[{
            "weight_kg":95.0,
            "reps":8,
            "rest_time_seconds":180,
            "rpe":9
          }]
        }],
        "notes":"curlでテスト！"
      }
    },
    "id":3
  }' | jq '.'

echo
echo "✅ curlテスト完了！"
