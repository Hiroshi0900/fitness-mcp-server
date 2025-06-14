package main

import (
	"database/sql"
	"fitness-mcp-server/internal/application/command/handler"
	command_usecase "fitness-mcp-server/internal/application/command/usecase"
	query_handler "fitness-mcp-server/internal/application/query/handler"
	query_usecase "fitness-mcp-server/internal/application/query/usecase"
	"fitness-mcp-server/internal/config"
	sqlite_query "fitness-mcp-server/internal/infrastructure/query/sqlite"
	"fitness-mcp-server/internal/infrastructure/repository/sqlite"
	"fitness-mcp-server/internal/interface/mcp/tool"
	"fitness-mcp-server/internal/interface/repository"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/server"
	_ "modernc.org/sqlite"
)

func main() {
	// 設定の初期化
	cfg := config.NewConfig()

	// データベースディレクトリを作成
	if err := cfg.EnsureDatabaseDir(); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// 依存関係の初期化
	dependencies, err := initializeDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// MCPサーバの作成
	s := server.NewMCPServer(
		cfg.MCP.Name,
		cfg.MCP.Version,
		server.WithToolCapabilities(false),
	)

	// ツールの登録
	if err := registerAllTools(s, dependencies); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// サーバの起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Dependencies はアプリケーションの依存関係を表します
type Dependencies struct {
	CommandHandler *handler.StrengthCommandHandler
	QueryHandler   *query_handler.StrengthQueryHandler
}

// initializeDependencies は依存関係を初期化します
func initializeDependencies(cfg *config.Config) (*Dependencies, error) {
	// リポジトリを初期化
	repo, err := initializeStrengthRepository(cfg.Database.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize strength repository: %w", err)
	}

	// クエリサービスを初期化
	queryService, err := initializeStrengthQueryService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize strength query service: %w", err)
	}

	// Command系の初期化
	commandUsecase := command_usecase.NewStrengthTrainingUsecase(repo)
	commandHandler := handler.NewStrengthCommandHandler(commandUsecase)

	// Query系の初期化
	queryUsecase := query_usecase.NewStrengthQueryUsecase(queryService)
	personalRecordsUsecase := query_usecase.NewPersonalRecordsUsecase(queryService)
	queryHandler := query_handler.NewStrengthQueryHandler(queryUsecase, personalRecordsUsecase)

	return &Dependencies{
		CommandHandler: commandHandler,
		QueryHandler:   queryHandler,
	}, nil
}

// registerAllTools はすべてのツールを登録します
func registerAllTools(s *server.MCPServer, deps *Dependencies) error {
	// トレーニング記録ツール
	trainingTool := tool.NewTrainingToolHandler(deps.CommandHandler)
	if err := trainingTool.Register(s); err != nil {
		return fmt.Errorf("failed to register training tool: %w", err)
	}

	// 期間指定クエリツール
	queryTool := tool.NewQueryToolHandler(deps.QueryHandler)
	if err := queryTool.Register(s); err != nil {
		return fmt.Errorf("failed to register query tool: %w", err)
	}

	// 個人記録ツール
	recordTool := tool.NewRecordToolHandler(deps.QueryHandler)
	if err := recordTool.Register(s); err != nil {
		return fmt.Errorf("failed to register record tool: %w", err)
	}

	return nil
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

// initializeStrengthQueryService はStrengthQueryServiceを初期化します
func initializeStrengthQueryService(cfg *config.Config) (*sqlite_query.StrengthQueryService, error) {
	// データベース接続を開く
	db, err := sql.Open("sqlite", cfg.Database.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQLiteの設定
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Hour)

	// SQLiteクエリサービスを作成
	queryService := sqlite_query.NewStrengthQueryService(db)

	log.Printf("Initialized SQLite query service at: %s", cfg.Database.SQLitePath)
	return queryService, nil
}
