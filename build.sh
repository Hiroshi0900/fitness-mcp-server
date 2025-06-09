#!/bin/bash

set -e

echo "Building fitness-mcp-server..."

# プロジェクトのルートディレクトリに移動
cd "$(dirname "$0")"

# 必要なディレクトリを作成
echo "Creating necessary directories..."
mkdir -p ~/.local/share/fitness-mcp

# Go プログラムをビルド
echo "Building Go binary..."
go build -o mcp-server ./cmd/server

echo "Build completed successfully!"
echo "Binary: $(pwd)/mcp-server"
echo "Data directory: ~/.local/share/fitness-mcp"

# 実行権限を確認
chmod +x mcp-server

echo "Setup complete!"
