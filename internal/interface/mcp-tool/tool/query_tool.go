package tool

import (
	"context"
	query_dto "fitness-mcp-server/internal/application/query/dto"
	query_handler "fitness-mcp-server/internal/application/query/handler"
	"fitness-mcp-server/internal/interface/mcp-tool/converter"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// QueryToolHandler は期間指定クエリツールを管理します
type QueryToolHandler struct {
	queryHandler *query_handler.StrengthQueryHandler
}

// NewQueryToolHandler は新しいQueryToolHandlerを作成します
func NewQueryToolHandler(queryHandler *query_handler.StrengthQueryHandler) *QueryToolHandler {
	return &QueryToolHandler{
		queryHandler: queryHandler,
	}
}

// Register は期間指定トレーニング取得ツールを登録します
func (h *QueryToolHandler) Register(s *server.MCPServer) error {
	tool := mcp.NewTool(
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

	toolHandler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return h.handleGetTrainingsByDateRange(ctx, req)
	}

	s.AddTool(tool, toolHandler)
	return nil
}

// handleGetTrainingsByDateRange は期間指定トレーニング取得処理を行います
func (h *QueryToolHandler) handleGetTrainingsByDateRange(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		response, err := h.queryHandler.GetTrainingsByDateRange(query)
		if err != nil {
			errorCh <- fmt.Errorf("トレーニング取得に失敗しました: %w", err)
			return
		}

		// レスポンスの整形
		result := converter.FormatQueryResponse(response)
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
