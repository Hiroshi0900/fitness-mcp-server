package command

import (
	"fitness-mcp-server/internal/application/dto"
	"fitness-mcp-server/internal/application/usecase"
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

//// RecordBIG3Session はBIG3セッションを記録する便利メソッド
//func (h *StrengthCommandHandler) RecordBIG3Session(
//	benchPressSets []SetCommand,
//	squatSets []SetCommand,
//	deadliftSets []SetCommand,
//	notes string,
//) (*RecordTrainingResult, error) {
//	return h.usecase.RecordBIG3Session(benchPressSets, squatSets, deadliftSets, notes)
//}
//
//// RecordQuickBenchPress はベンチプレスの簡単記録メソッド
//func (h *StrengthCommandHandler) RecordQuickBenchPress(weightKg float64, reps int, sets int, notes string) (*RecordTrainingResult, error) {
//	return h.usecase.RecordQuickBenchPress(weightKg, reps, sets, notes)
//}
