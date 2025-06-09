package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fitness-mcp-server/internal/application/command"
	"fitness-mcp-server/internal/application/usecase"
	"fitness-mcp-server/internal/config"
	"fitness-mcp-server/internal/infrastructure/mcp"
	"fitness-mcp-server/internal/infrastructure/repository/sqlite"
	"fitness-mcp-server/internal/interface/repository"
)

func main() {
	// 設定を読み込み
	cfg := config.NewConfig()

	// データベースディレクトリを作成
	if err := cfg.EnsureDatabaseDir(); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// リポジトリを初期化
	strengthRepo, err := initializeStrengthRepository(cfg.Database.SQLitePath)
	if err != nil {
		log.Fatalf("Failed to initialize strength repository: %v", err)
	}

	// ユースケースを初期化
	strengthUsecase := usecase.NewStrengthTrainingUsecase(strengthRepo)

	// コマンドハンドラーを初期化
	strengthHandler := command.NewStrengthCommandHandler(strengthUsecase)

	// MCPサーバーを作成
	mcpServer := mcp.NewMCPServer(
		cfg.MCP.Name,
		cfg.MCP.Version,
		cfg.MCP.Description,
	)

	// 筋トレツールを登録
	mcp.RegisterStrengthTools(mcpServer, strengthHandler)

	log.Printf("HTTP MCP Server starting: %s v%s", cfg.MCP.Name, cfg.MCP.Version)
	log.Printf("Database: %s", cfg.Database.SQLitePath)
	log.Printf("Registered tools: %d", len(mcpServer.GetTools()))
	log.Printf("Listening on http://localhost:8080")

	// HTTPハンドラーを設定
	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		// リクエストボディを読み取り
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// MCPサーバーで処理
		response := mcpServer.HandleRequest(body)

		// レスポンスを返す
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	// ヘルスチェック用エンドポイント
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"server": cfg.MCP.Name,
			"version": cfg.MCP.Version,
		})
	})

	// サーバー起動
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// initializeStrengthRepository はStrengthRepositoryを初期化します
func initializeStrengthRepository(dbPath string) (repository.StrengthTrainingRepository, error) {
	// SQLiteリポジトリを作成
	repo, err := sqlite.NewStrengthRepository(dbPath)
	if err != nil {
		return nil, err
	}

	// データベースの初期化（テーブル作成）
	if err := repo.Initialize(); err != nil {
		return nil, err
	}

	log.Printf("Initialized SQLite repository at: %s", dbPath)
	return repo, nil
}
