package handler

import (
	"fitness-mcp-server/internal/application/command/dto"
	"fitness-mcp-server/internal/application/command/usecase"
)

// =============================================================================
// 筋トレコマンドハンドラー - ビジネスロジックの実行
// =============================================================================

// StrengthCommandHandler は筋トレに関するコマンドを処理するハンドラー
type StrengthCommandHandler struct {
	usecase usecase.StrengthTrainingUsecase
}

// NewStrengthCommandHandler は新しいStrengthCommandHandlerを作成します
func NewStrengthCommandHandler(usecase usecase.StrengthTrainingUsecase) *StrengthCommandHandler {
	return &StrengthCommandHandler{
		usecase: usecase,
	}
}

// RecordTraining は筋トレセッションを記録します
func (h *StrengthCommandHandler) RecordTraining(cmd dto.RecordTrainingCommand) (*dto.RecordTrainingResult, error) {
	return h.usecase.RecordTraining(cmd)
}

// UpdateTraining は筋トレセッションを更新します
func (h *StrengthCommandHandler) UpdateTraining(cmd dto.UpdateTrainingCommand) (*dto.UpdateTrainingResult, error) {
	return h.usecase.UpdateTraining(cmd)
}

// DeleteTraining は筋トレセッションを削除します
func (h *StrengthCommandHandler) DeleteTraining(cmd dto.DeleteTrainingCommand) (*dto.DeleteTrainingResult, error) {
	return h.usecase.DeleteTraining(cmd)
}
