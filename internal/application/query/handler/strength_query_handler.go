package handler

import (
	"fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/application/query/usecase"
)

// StrengthQueryHandler は筋トレデータの読み取り系ハンドラー
type StrengthQueryHandler struct {
	usecase           usecase.StrengthQueryUsecase
	personalRecordsUC usecase.PersonalRecordsUsecase
}

// NewStrengthQueryHandler は新しいStrengthQueryHandlerを作成します
func NewStrengthQueryHandler(
	usecase usecase.StrengthQueryUsecase,
	personalRecordsUC usecase.PersonalRecordsUsecase,
) *StrengthQueryHandler {
	return &StrengthQueryHandler{
		usecase:           usecase,
		personalRecordsUC: personalRecordsUC,
	}
}

// GetTrainingsByDateRange は指定した期間のトレーニングセッションを取得します
func (h *StrengthQueryHandler) GetTrainingsByDateRange(query dto.GetTrainingsByDateRangeQuery) (*dto.GetTrainingsByDateRangeResponse, error) {
	return h.usecase.GetTrainingsByDateRange(query)
}

// GetPersonalRecords は個人記録を取得します
func (h *StrengthQueryHandler) GetPersonalRecords(query dto.GetPersonalRecordsQuery) (*dto.GetPersonalRecordsResponse, error) {
	return h.personalRecordsUC.GetPersonalRecords(query)
}
