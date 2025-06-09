#!/bin/bash

echo "Testing fitness-mcp-server startup..."

# 環境変数設定
export MCP_DATA_DIR="/Users/sakemihiroshi/.local/share/fitness-mcp"

# MCPサーバーを起動してテスト
cd "$(dirname "$0")"

echo "Starting MCP server with data directory: $MCP_DATA_DIR"

# JSON-RPC initializeメッセージを送信
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | ./mcp-server

echo "Test completed."
