#!/bin/bash

# fitness-mcp-server CLI実行サンプル
# MCPサーバーとJSON-RPC形式で通信します

SERVER_BINARY="./mcp"
DATE=$(date +%Y-%m-%d)

echo "=== Fitness MCP Server CLI Examples ==="

# サーバーの起動確認
if [ ! -f "$SERVER_BINARY" ]; then
    echo "Error: MCPサーバーバイナリが見つかりません。'go build ./cmd/mcp'を実行してください。"
    exit 1
fi

# 1. トレーニング記録の追加
echo "1. トレーニング記録を追加..."
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "record_training",
    "arguments": {
      "date": "'$DATE'",
      "exercises": [
        {
          "name": "ベンチプレス",
          "category": "Compound",
          "sets": [
            {
              "weight_kg": 80.0,
              "reps": 10,
              "rest_time_seconds": 120,
              "rpe": 7
            },
            {
              "weight_kg": 80.0,
              "reps": 8,
              "rest_time_seconds": 120,
              "rpe": 8
            }
          ]
        },
        {
          "name": "ダンベルカール",
          "category": "Isolation",
          "sets": [
            {
              "weight_kg": 15.0,
              "reps": 12,
              "rest_time_seconds": 60,
              "rpe": 6
            }
          ]
        }
      ],
      "notes": "調子良好、フォームを意識"
    }
  }
}' | $SERVER_BINARY

echo -e "\n---\n"

# 2. 期間指定でトレーニング履歴を取得
echo "2. 今月のトレーニング履歴を取得..."
START_DATE=$(date +%Y-%m-01)
END_DATE=$(date +%Y-%m-%d)

echo '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get_trainings_by_date_range",
    "arguments": {
      "start_date": "'$START_DATE'",
      "end_date": "'$END_DATE'"
    }
  }
}' | $SERVER_BINARY

echo -e "\n---\n"

# 3. 個人記録を取得
echo "3. 全ての個人記録を取得..."
echo '{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_personal_records",
    "arguments": {}
  }
}' | $SERVER_BINARY

echo -e "\n---\n"

# 4. 特定のエクササイズの個人記録を取得
echo "4. ベンチプレスの個人記録を取得..."
echo '{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "get_personal_records",
    "arguments": {
      "exercise_name": "ベンチプレス"
    }
  }
}' | $SERVER_BINARY

echo -e "\n=== 完了 ==="
