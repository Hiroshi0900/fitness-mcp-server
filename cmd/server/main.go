package main

import (
	"bufio"
	"log"
	"os"

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

	log.Printf("MCP Server starting: %s v%s", cfg.MCP.Name, cfg.MCP.Version)
	log.Printf("Database: %s", cfg.Database.SQLitePath)
	log.Printf("Registered tools: %d", len(mcpServer.GetTools()))

	// stdin/stdout でJSON-RPC通信を開始
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// リクエストを処理してレスポンスを出力
		response := mcpServer.HandleRequest(line)
		if response != nil { // nilの場合は出力しない（通知メッセージの場合）
			os.Stdout.Write(response)
			os.Stdout.Write([]byte("\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
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
