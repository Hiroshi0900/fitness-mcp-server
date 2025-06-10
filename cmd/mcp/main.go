package main

import (
	"context"
	"fitness-mcp-server/internal/application/command"
	"fitness-mcp-server/internal/application/dto"
	query_dto "fitness-mcp-server/internal/application/query/dto"
	query_handler "fitness-mcp-server/internal/application/query/handler"
	query_usecase "fitness-mcp-server/internal/application/query/usecase"
	"fitness-mcp-server/internal/application/usecase"
	"fitness-mcp-server/internal/config"
	"fitness-mcp-server/internal/infrastructure/repository/sqlite"
	"fitness-mcp-server/internal/interface/repository"
	"fmt"
	"log"
	"time"

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

	// Command系の初期化
	commandUsecase := usecase.NewStrengthTrainingUsecase(repo)
	commandHandler := command.NewStrengthCommandHandler(commandUsecase)

	// Query系の初期化
	queryUsecase := query_usecase.NewStrengthQueryUsecase(repo)
	queryHandler := query_handler.NewStrengthQueryHandler(queryUsecase)

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

		result, err := commandHandler.RecordTraining(cmd)
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

	// クエリツールの追加
	queryTool := mcp.NewTool(
		"get_trainings_by_date_range",
		mcp.WithDescription("指定した期間のトレーニングセッションを取得する"),
		mcp.WithString("start_date",
			mcp.Required(),
			mcp.Description("検索開始日（YYYY-MM-DD形式）"),
		),
		mcp.WithString("end_date",
			mcp.Required(),
			mcp.Description("検索終了日（YYYY-MM-DD形式）"),
		),
	)

	queryToolHandler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// タイムアウト設定（30秒）
		timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Goroutineで処理を実行
		resultCh := make(chan *mcp.CallToolResult, 1)
		errorCh := make(chan error, 1)

		go func() {
			// パラメータの取得
			startDateStr, err := req.RequireString("start_date")
			if err != nil {
				errorCh <- fmt.Errorf("start_date パラメータが必要です: %w", err)
				return
			}

			endDateStr, err := req.RequireString("end_date")
			if err != nil {
				errorCh <- fmt.Errorf("end_date パラメータが必要です: %w", err)
				return
			}

			// 日付のパース
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				errorCh <- fmt.Errorf("start_date の形式が不正です: %w", err)
				return
			}

			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				errorCh <- fmt.Errorf("end_date の形式が不正です: %w", err)
				return
			}

			// クエリの実行
			query := query_dto.GetTrainingsByDateRangeQuery{
				StartDate: startDate,
				EndDate:   endDate,
			}

			response, err := queryHandler.GetTrainingsByDateRange(query)
			if err != nil {
				errorCh <- fmt.Errorf("トレーニング取得に失敗しました: %w", err)
				return
			}

			// レスポンスの整形
			result := formatQueryResponse(response)
			resultCh <- mcp.NewToolResultText(result)
		}()

		// タイムアウトまたは結果を待機
		select {
		case <-timeoutCtx.Done():
			return mcp.NewToolResultError("リクエストがタイムアウトしました（30秒）"), nil
		case err := <-errorCh:
			return mcp.NewToolResultError(err.Error()), nil
		case result := <-resultCh:
			return result, nil
		}
	}

	s.AddTool(queryTool, queryToolHandler)

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

// formatQueryResponse はクエリレスポンスを見やすい形式にフォーマットします（簡略版）
func formatQueryResponse(response *query_dto.GetTrainingsByDateRangeResponse) string {
	if response.Count == 0 {
		return fmt.Sprintf("📊 **期間: %s**\n\n❌ この期間にトレーニング記録は見つかりませんでした。", response.Period)
	}

	result := fmt.Sprintf("📊 **期間: %s**\n\n🏋️ **トレーニング記録: %d件**\n\n", response.Period, response.Count)

	for i, training := range response.Trainings {
		result += fmt.Sprintf("**%d. %s (%s)**\n",
			i+1,
			training.Date.Format("2006-01-02"),
			training.Date.Weekday())

		if training.Notes != "" {
			result += fmt.Sprintf("📝 メモ: %s\n", training.Notes)
		}

		result += fmt.Sprintf("📈 概要: %d種目, %dセット, %.1fkg総ボリューム\n",
			training.Summary.TotalExercises,
			training.Summary.TotalSets,
			training.Summary.TotalVolume)

		// エクササイズの概要のみ（詳細は省略）
		for _, exercise := range training.Exercises {
			result += fmt.Sprintf("  • %s (%s): %d sets\n",
				exercise.Name, exercise.Category, len(exercise.Sets))
		}
		result += "\n"
	}

	return result
}
