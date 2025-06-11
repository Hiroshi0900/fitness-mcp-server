#!/bin/bash

echo "=== 記録 ==="
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call","params": {"name": "record_training", "arguments": {"date": "2025-06-08", "exercises": [{"name": "ベンチプレス", "category": "Compound", "sets": [{"weight_kg": 90, "reps": 10, "rest_time_seconds": 180}]}]}}}' | ./mcp-server
