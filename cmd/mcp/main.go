package main

import (
	"context"
	"fitness-mcp-server/internal/application/command"
	"fitness-mcp-server/internal/application/dto"
	"fitness-mcp-server/internal/application/usecase"
	"fitness-mcp-server/internal/config"
	"fitness-mcp-server/internal/infrastructure/repository/sqlite"
	"fitness-mcp-server/internal/interface/repository"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 依存の初期化
	cfg := config.NewConfig()
	// データベースディレクトリを作成
	if err := cfg.EnsureDatabaseDir(); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// リポジトリを初期化
	repo, err := initializeStrengthRepository(cfg.Database.SQLitePath)
	if err != nil {
		log.Fatalf("Failed to initialize strength repository: %v", err)
	}
	usecase := usecase.NewStrengthTrainingUsecase(repo)
	handler := command.NewStrengthCommandHandler(usecase)

	// ToolHandlerFuncのラップ
	toolHandlerFunc := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 必須パラメータ取得
		fitness, err := req.RequireString("fitness")
		if err != nil {
			return mcp.NewToolResultError("fitnessパラメータが必要です: " + err.Error()), nil
		}

		// RecordTrainingCommandのNotesにfitnessを入れる（最低限の例）
		cmd := dto.RecordTrainingCommand{
			Notes: fitness,
			// DateやExercisesは本来必須だが、ここでは省略（本番では要対応）
		}

		result, err := handler.RecordTraining(cmd)
		if err != nil {
			return mcp.NewToolResultError("記録に失敗しました: " + err.Error()), nil
		}

		// 結果をテキストで返す
		return mcp.NewToolResultText(
			fmt.Sprintf("記録完了: TrainingID=%v, メッセージ=%v", result.TrainingID, result.Message),
		), nil
	}

	// サーバの起動
	s := server.NewMCPServer(
		"筋トレ記録サーバ",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// ツールの登録
	tool := mcp.NewTool(
		"筋トレ記録ツール",
		mcp.WithDescription("筋トレの記録を管理するツール"),
		mcp.WithString("fitness",
			mcp.Required(),
			mcp.Description("筋トレの種類?"),
		),
	)

	// ツールをサーバに登録
	s.AddTool(tool, toolHandlerFunc)
	
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
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
