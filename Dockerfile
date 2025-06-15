# Build stage
FROM golang:1.24.4-alpine AS builder

# 必要なパッケージのインストール
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 作業ディレクトリの設定
WORKDIR /app

# Go modulesのコピーと依存関係の取得
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# バイナリのビルド
# CGO_ENABLED=1 でSQLiteのサポートを有効にする
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o fitness-mcp-server ./cmd/mcp

# Runtime stage
FROM alpine:latest

# 必要なパッケージのインストール
RUN apk --no-cache add ca-certificates sqlite

# 非rootユーザーの作成
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 作業ディレクトリの設定
WORKDIR /app

# ビルドしたバイナリのコピー
COPY --from=builder /app/fitness-mcp-server .

# データディレクトリの作成と権限設定
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

# 非rootユーザーに切り替え
USER appuser

# 環境変数の設定
ENV MCP_DATA_DIR=/app/data

# ヘルスチェック（オプション）
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD test -f /app/data/fitness.db || exit 1

# エントリーポイント
CMD ["./fitness-mcp-server"]
