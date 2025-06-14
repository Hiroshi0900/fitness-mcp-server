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

// RecordToolHandler は個人記録取得ツールを管理します
type RecordToolHandler struct {
	queryHandler *query_handler.StrengthQueryHandler
}

// NewRecordToolHandler は新しいRecordToolHandlerを作成します
func NewRecordToolHandler(queryHandler *query_handler.StrengthQueryHandler) *RecordToolHandler {
	return &RecordToolHandler{
		queryHandler: queryHandler,
	}
}

// Register は個人記録取得ツールを登録します
func (h *RecordToolHandler) Register(s *server.MCPServer) error {
	tool := mcp.NewTool(
		"get_personal_records",
		mcp.WithDescription("個人記録（最大重量、最大レップ数、最大ボリューム等）を取得する"),
		mcp.WithString("exercise_name",
			mcp.Description("特定のエクササイズ名（省略可）。指定すると該当エクササイズの記録のみを取得します。"),
		),
	)

	toolHandler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return h.handleGetPersonalRecords(ctx, req)
	}

	s.AddTool(tool, toolHandler)
	return nil
}

// handleGetPersonalRecords は個人記録取得処理を行います
func (h *RecordToolHandler) handleGetPersonalRecords(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// タイムアウト設定（30秒）
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Goroutineで処理を実行
	resultCh := make(chan *mcp.CallToolResult, 1)
	errorCh := make(chan error, 1)

	go func() {
		// パラメータの取得（オプション）
		var exerciseName *string
		if paramsMap, ok := req.Params.Arguments.(map[string]interface{}); ok {
			if name, exists := paramsMap["exercise_name"]; exists {
				if nameStr, ok := name.(string); ok && nameStr != "" {
					exerciseName = &nameStr
				}
			}
		}

		// クエリの実行
		query := query_dto.GetPersonalRecordsQuery{
			ExerciseName: exerciseName,
		}

		response, err := h.queryHandler.GetPersonalRecords(query)
		if err != nil {
			errorCh <- fmt.Errorf("個人記録取得に失敗しました: %w", err)
			return
		}

		// レスポンスの整形
		result := converter.FormatPersonalRecordsResponse(response)
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
