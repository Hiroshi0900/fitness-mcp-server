#!/bin/bash

# フィットネスMCPサーバーの実行例

echo "🏋️ Fitness MCP Server - 実行例"

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# バイナリをビルド
echo "📦 バイナリをビルド中..."
go build -o mcp ./cmd/mcp

if [ $? -ne 0 ]; then
    echo "❌ ビルドに失敗しました"
    exit 1
fi

echo "✅ ビルド完了"

# 実行権限を付与
chmod +x scripts/cli_examples.sh
chmod +x scripts/fitness_cli.py

echo ""
echo "🎯 利用可能なコマンド例:"
echo ""

echo "1. シェルスクリプトでサンプル実行:"
echo "   ./scripts/cli_examples.sh"
echo ""

echo "2. Python CLIツール:"
echo "   # インタラクティブな記録"
echo "   python3 scripts/fitness_cli.py record"
echo ""
echo "   # 今月の履歴表示"
echo "   python3 scripts/fitness_cli.py history --start $(date +%Y-%m-01) --end $(date +%Y-%m-%d)"
echo ""
echo "   # 個人記録表示"
echo "   python3 scripts/fitness_cli.py records"
echo ""
echo "   # 特定エクササイズの記録"
echo "   python3 scripts/fitness_cli.py records --exercise ベンチプレス"
echo ""
echo "   # クイック記録"
echo "   python3 scripts/fitness_cli.py quick"
echo ""

echo "3. 直接JSON-RPC実行:"
echo '   echo '"'"'{
     "jsonrpc": "2.0",
     "id": 1,
     "method": "tools/call",
     "params": {
       "name": "get_personal_records",
       "arguments": {}
     }
   }'"'"' | ./mcp'
echo ""

echo "🚀 実行してみましょう！"
