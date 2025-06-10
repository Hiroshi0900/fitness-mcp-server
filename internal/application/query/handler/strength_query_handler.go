package handler

import (
	"fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/application/query/usecase"
)

// StrengthQueryHandler は筋トレデータの読み取り系ハンドラー
type StrengthQueryHandler struct {
	usecase usecase.StrengthQueryUsecase
}

// NewStrengthQueryHandler は新しいStrengthQueryHandlerを作成します
func NewStrengthQueryHandler(usecase usecase.StrengthQueryUsecase) *StrengthQueryHandler {
	return &StrengthQueryHandler{
		usecase: usecase,
	}
}

// GetTrainingsByDateRange は指定した期間のトレーニングセッションを取得します
func (h *StrengthQueryHandler) GetTrainingsByDateRange(query dto.GetTrainingsByDateRangeQuery) (*dto.GetTrainingsByDateRangeResponse, error) {
	return h.usecase.GetTrainingsByDateRange(query)
}