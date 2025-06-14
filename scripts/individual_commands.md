# Individual Command Examples

## 1. トレーニング記録

```bash
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "record_training",
    "arguments": {
      "date": "2024-06-14",
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
              "weight_kg": 82.5,
              "reps": 8,
              "rest_time_seconds": 120,
              "rpe": 8
            },
            {
              "weight_kg": 85.0,
              "reps": 6,
              "rest_time_seconds": 180,
              "rpe": 9
            }
          ]
        },
        {
          "name": "スクワット",
          "category": "Compound",
          "sets": [
            {
              "weight_kg": 100.0,
              "reps": 12,
              "rest_time_seconds": 120,
              "rpe": 6
            },
            {
              "weight_kg": 110.0,
              "reps": 10,
              "rest_time_seconds": 120,
              "rpe": 7
            }
          ]
        }
      ],
      "notes": "フォームを意識、RPE高め"
    }
  }
}' | ./mcp
```

## 2. 期間指定でトレーニング履歴取得

```bash
echo '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get_trainings_by_date_range",
    "arguments": {
      "start_date": "2024-06-01",
      "end_date": "2024-06-14"
    }
  }
}' | ./mcp
```

## 3. 全個人記録取得

```bash
echo '{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_personal_records",
    "arguments": {}
  }
}' | ./mcp
```

## 4. 特定エクササイズの個人記録

```bash
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
}' | ./mcp
```

## 5. ワンライナーでクイック記録

```bash
# 今日のベンチプレス記録
DATE=$(date +%Y-%m-%d)
echo "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"record_training\",\"arguments\":{\"date\":\"$DATE\",\"exercises\":[{\"name\":\"ベンチプレス\",\"category\":\"Compound\",\"sets\":[{\"weight_kg\":80.0,\"reps\":10,\"rest_time_seconds\":120}]}]}}}" | ./mcp
```

## 6. 環境変数を使った設定

```bash
# カスタムデータディレクトリ
export MCP_DATA_DIR="./custom_data"
export MCP_SERVER_NAME="My Fitness Server"
export DB_MAX_OPEN_CONNS=20

echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_personal_records",
    "arguments": {}
  }
}' | ./mcp
```

## 7. 複数の記録を連続実行

```bash
#!/bin/bash

# 週間ワークアウトの一括記録
EXERCISES=(
  '{"name":"ベンチプレス","category":"Compound","sets":[{"weight_kg":80,"reps":10,"rest_time_seconds":120}]}'
  '{"name":"スクワット","category":"Compound","sets":[{"weight_kg":100,"reps":12,"rest_time_seconds":120}]}'
  '{"name":"デッドリフト","category":"Compound","sets":[{"weight_kg":120,"reps":8,"rest_time_seconds":180}]}'
)

for i in "${!EXERCISES[@]}"; do
  DATE=$(date -d "+$i days" +%Y-%m-%d)
  echo "{\"jsonrpc\":\"2.0\",\"id\":$((i+1)),\"method\":\"tools/call\",\"params\":{\"name\":\"record_training\",\"arguments\":{\"date\":\"$DATE\",\"exercises\":[${EXERCISES[$i]}]}}}" | ./mcp
  echo "---"
done
```

## 実行方法

1. プロジェクトルートで実行:
```bash
cd fitness-mcp-server
go build -o mcp ./cmd/mcp
```

2. 上記のコマンドをコピー&ペーストして実行

3. または準備されたスクリプトを使用:
```bash
./scripts/run_examples.sh
python3 scripts/fitness_cli.py record
```
